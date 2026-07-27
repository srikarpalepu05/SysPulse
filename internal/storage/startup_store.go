package storage

import (
	"database/sql"
	"fmt"
	"time"
)

type StartupRecord struct {
	Name      string
	Command   string
	Source    string
	Location  string
	Timestamp time.Time
}

func EnsureStartupTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS startup_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp TEXT NOT NULL,
			name TEXT NOT NULL,
			command TEXT NOT NULL,
			source TEXT NOT NULL,
			location TEXT NOT NULL
		);`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("create startup_items table: %w", err)
	}

	return nil
}

func SaveStartupEntry(db *sql.DB, entry StartupRecord) error {
	query := `
		INSERT INTO startup_items (timestamp, name, command, source, location)
		VALUES (?, ?, ?, ?, ?)`

	_, err := db.Exec(query, entry.Timestamp.Format(time.RFC3339), entry.Name, entry.Command, entry.Source, entry.Location)
	if err != nil {
		return fmt.Errorf("insert startup entry: %w", err)
	}

	return nil
}
