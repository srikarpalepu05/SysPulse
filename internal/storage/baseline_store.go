package storage

import (
	"database/sql"
	"fmt"
	"time"
)

type BaselineSnapshot struct {
	GeneratedAt  time.Time
	StartupCount int
	ProcessCount int
	Summary      string
}

func EnsureBaselineTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS baseline_snapshots (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			generated_at TEXT NOT NULL,
			startup_count INTEGER NOT NULL,
			process_count INTEGER NOT NULL,
			summary TEXT NOT NULL
		);`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("create baseline_snapshots table: %w", err)
	}

	return nil
}

func SaveBaseline(db *sql.DB, snapshot BaselineSnapshot) error {
	query := `
		INSERT INTO baseline_snapshots (generated_at, startup_count, process_count, summary)
		VALUES (?, ?, ?, ?)`

	_, err := db.Exec(query, snapshot.GeneratedAt.Format(time.RFC3339), snapshot.StartupCount, snapshot.ProcessCount, snapshot.Summary)
	if err != nil {
		return fmt.Errorf("insert baseline snapshot: %w", err)
	}

	return nil
}

func GetLatestBaseline(db *sql.DB) (BaselineSnapshot, error) {
	var snapshot BaselineSnapshot
	var generatedAt string

	query := `
		SELECT generated_at, startup_count, process_count, summary
		FROM baseline_snapshots
		ORDER BY id DESC
		LIMIT 1`

	row := db.QueryRow(query)
	if err := row.Scan(&generatedAt, &snapshot.StartupCount, &snapshot.ProcessCount, &snapshot.Summary); err != nil {
		return BaselineSnapshot{}, err
	}

	parsed, err := time.Parse(time.RFC3339, generatedAt)
	if err != nil {
		return BaselineSnapshot{}, fmt.Errorf("parse baseline timestamp: %w", err)
	}

	snapshot.GeneratedAt = parsed
	return snapshot, nil
}
