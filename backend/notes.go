package main

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

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
	_, err := db.Exec(
		"INSERT INTO notes (name, message, ip, created_at) VALUE (?, ?, ?, ?)",
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
