package main

import (
	"log"
	"net/http"
)

func main() {
	db, err := openDB("notes.db")
	if err != nil {
		log.Fatal("opening database: ", err)
	}
	defer db.Close()

	httpAddr := "127.0.0.1:3000"
	sshAddr := "127.0.0.1:2222"

	// http.Handle("GET /", http.FileServer(http.Dir("../site")))
	http.HandleFunc("GET /api/projects", handleProjects)
	http.HandleFunc("GET /api/notes", handleNotes(db))

	sshServer, err := newSSHServer(sshAddr, db)
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
