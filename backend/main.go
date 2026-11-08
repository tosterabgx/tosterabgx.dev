package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	addr := "127.0.0.1:3000"

	http.HandleFunc("GET /api/projects", handleProjects)

	fmt.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
