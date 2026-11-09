package main

import (
	"log"
	"net/http"
	"os"

	"charm.land/ssh"
)

func main() {
	httpAddr := envOr("HTTP_ADDR", "127.0.0.1:3000")
	sshAddr := envOr("SSH_ADDR", "127.0.0.1:2222")
	dbPath := envOr("DB_PATH", "notes.db")

	db, err := openDB(dbPath)
	if err != nil {
		log.Fatal("opening database: ", err)
	}
	defer db.Close()

	http.Handle("GET /", http.FileServer(http.Dir("../site")))
	http.HandleFunc("GET /api/projects", handleProjects)
	http.HandleFunc("GET /api/notes", handleNotes(db))

	sshServer, err := newSSHServer(sshAddr, db, loadAdminKey())
	if err != nil {
		log.Fatal("creating ssh server: ", err)
	}

	go func() {
		log.Printf("ssh listening on %s\n", sshAddr)
		log.Fatal(sshServer.ListenAndServe())
	}()

	log.Printf("http listening on %s\n", httpAddr)
	log.Fatal(http.ListenAndServe(httpAddr, nil))
}

func envOr(key, d string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return d
}

func loadAdminKey() ssh.PublicKey {
	s := os.Getenv("ADMIN_KEY")
	if s == "" {
		log.Println("ADMIN_KEY not set, admin mode disabled")
		return nil
	}
	key, _, _, _, err := ssh.ParseAuthorizedKey([]byte(s))
	if err != nil {
		log.Fatal("parsing ADMIN_KEY: ", err)
	}
	log.Println("admin mode enabled")
	return key
}
