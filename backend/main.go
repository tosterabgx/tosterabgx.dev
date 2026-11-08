package main

import (
	"log"
	"net/http"
)

func main() {
	httpAddr := "127.0.0.1:3000"
	sshAddr := "127.0.0.1:2222"

	http.HandleFunc("GET /api/projects", handleProjects)

	sshServer, err := newSSHServer(sshAddr)
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
