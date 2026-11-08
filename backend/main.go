package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Project struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

func handleProjects(w http.ResponseWriter, r *http.Request) {
	projects := []Project{{"Project1", "Description of project1", "url"}, {"Project2", "Description of project2", "url"}}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(projects)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
	}
}

var port uint16 = 3000

func main() {
	http.HandleFunc("GET /api/projects", handleProjects)
	fmt.Printf("listening on 127.0.0.1:%d", port)

	if err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), nil); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
	}
}
