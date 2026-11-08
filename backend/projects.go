package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type Project struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

type githubRepo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	HTMLURL     string `json:"html_url"`
	Language    string `json:"language"`
	Fork        bool   `json:"fork"`
}

const githubReposURL = "https://api.github.com/users/tosterabgx/repos?sort=pushed&per_page=100"

var client = &http.Client{Timeout: 10 * time.Second}

var (
	cachedProjects []Project
	cachedAt       time.Time
	cacheMu        sync.Mutex
)

const cacheTTL = 30 * time.Minute

func fetchProjects() ([]Project, error) {
	resp, err := client.Get(githubReposURL)
	if err != nil {
		return nil, fmt.Errorf("requesting github: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("github returned status %d", resp.StatusCode)
	}

	var repos []githubRepo
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, fmt.Errorf("decoding github response: %w", err)
	}

	var projects []Project
	for _, r := range repos {
		if r.Fork {
			continue
		}
		p := Project{
			Name:        r.Name,
			Description: r.Description,
			URL:         r.HTMLURL,
		}
		projects = append(projects, p)
	}

	return projects, nil
}

func getProjects() ([]Project, error) {
	cacheMu.Lock()
	defer cacheMu.Unlock()

	if cachedProjects != nil && time.Since(cachedAt) < cacheTTL {
		return cachedProjects, nil
	}

	projects, err := fetchProjects()
	if err != nil {
		if cachedProjects != nil {
			log.Println("github failed, stale cache:", err)
			return cachedProjects, nil
		}

		return nil, err
	}

	cachedProjects = projects
	cachedAt = time.Now()
	return projects, nil
}

func handleProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := getProjects()
	if err != nil {
		log.Println("fetching projects:", err)
		http.Error(w, "error fetching github projects", 502)
		return
	}

	if projects == nil {
		projects = []Project{}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(projects); err != nil {
		log.Println("error:", err)
		return
	}
}
