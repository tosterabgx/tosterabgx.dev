package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

const (
	maxNameLen = 24
	maxNoteLen = 100
)

const (
	rateLimitCount  = 3
	rateLimitWindow = 10 * time.Minute
)

var errRatedLimited = errors.New("too many notes from this address")

type Note struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	IP        string    `json:"-"`
}

func openDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", "file:"+path+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS notes (
			id			INTEGER PRIMARY KEY AUTOINCREMENT,
			name		TEXT	NOT NULL,
			message 	TEXT	NOT NULL,
			ip			TEXT	NOT NULL,
			created_at	INTEGER NOT NULL
		)
	`)
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func insertNote(db *sql.DB, name, message, ip string) error {
	name = sanitize(name, maxNameLen)
	message = sanitize(message, maxNoteLen)
	if name == "" || message == "" {
		return errors.New("empty name or message")
	}

	n, err := countRecentNotes(db, ip, time.Now().Add(-rateLimitWindow))
	if err != nil {
		return err
	}
	if n >= rateLimitCount {
		return errRatedLimited
	}

	_, err = db.Exec(
		"INSERT INTO notes (name, message, ip, created_at) VALUES (?, ?, ?, ?)",
		name, message, ip, time.Now().Unix(),
	)
	return err
}

func listNotes(db *sql.DB, limit int) ([]Note, error) {
	rows, err := db.Query(
		"SELECT id, name, message, created_at FROM notes ORDER BY id DESC LIMIT ?",
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := []Note{}
	for rows.Next() {
		var n Note
		var unix int64
		if err := rows.Scan(&n.ID, &n.Name, &n.Message, &unix); err != nil {
			return nil, err
		}
		n.CreatedAt = time.Unix(unix, 0).UTC()
		notes = append(notes, n)
	}
	return notes, rows.Err()
}

func handleNotes(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		notes, err := listNotes(db, 50)
		if err != nil {
			log.Println("listing notes:", err)
			http.Error(w, "could not load notes", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(notes); err != nil {
			log.Println("writing notes response:", err)
		}
	}
}

func countRecentNotes(db *sql.DB, ip string, since time.Time) (int, error) {
	var n int
	err := db.QueryRow(
		"SELECT COUNT(*) FROM notes WHERE ip = ? AND created_at > ?",
		ip, since.Unix(),
	).Scan(&n)

	return n, err
}
