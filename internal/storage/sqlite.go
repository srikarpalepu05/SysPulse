package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
	"syspulse/internal/collector"
	"syspulse/internal/rules"
)

func OpenDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	createSnapshots := `
		CREATE TABLE IF NOT EXISTS process_snapshots (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp TEXT NOT NULL,
			process_name TEXT NOT NULL,
			pid INTEGER NOT NULL,
			cpu_percent REAL NOT NULL,
			memory_mb REAL NOT NULL
		);`

	createAlerts := `
		CREATE TABLE IF NOT EXISTS alerts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp TEXT NOT NULL,
			process_name TEXT NOT NULL,
			pid INTEGER NOT NULL,
			severity TEXT NOT NULL,
			reason TEXT NOT NULL,
			cpu_percent REAL NOT NULL,
			memory_mb REAL NOT NULL
		);`

	if _, err = db.Exec(createSnapshots); err != nil {
		return nil, err
	}

	if _, err = db.Exec(createAlerts); err != nil {
		return nil, err
	}

	return db, nil
}

func SaveSnapshot(db *sql.DB, snapshot collector.ProcessSnapshot) error {
	query := `
		INSERT INTO process_snapshots (timestamp, process_name, pid, cpu_percent, memory_mb)
		VALUES (?, ?, ?, ?, ?)`

	_, err := db.Exec(query, snapshot.Timestamp.Format(time.RFC3339), snapshot.Name, snapshot.PID, snapshot.CPU, snapshot.MemoryMB)
	if err != nil {
		return fmt.Errorf("insert snapshot: %w", err)
	}

	return nil
}

func SaveAlert(db *sql.DB, alert rules.Alert) error {
	query := `
		INSERT INTO alerts (timestamp, process_name, pid, severity, reason, cpu_percent, memory_mb)
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := db.Exec(query, time.Now().Format(time.RFC3339), alert.ProcessName, alert.PID, alert.Severity, alert.Reason, alert.CPU, alert.MemoryMB)
	if err != nil {
		return fmt.Errorf("insert alert: %w", err)
	}

	return nil
}

func ListAlerts(db *sql.DB) ([]rules.Alert, error) {
	rows, err := db.Query(`SELECT process_name, pid, cpu_percent, memory_mb, severity, reason FROM alerts ORDER BY id DESC LIMIT 20`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	alerts := make([]rules.Alert, 0)
	for rows.Next() {
		var alert rules.Alert
		if err := rows.Scan(&alert.ProcessName, &alert.PID, &alert.CPU, &alert.MemoryMB, &alert.Severity, &alert.Reason); err != nil {
			return nil, err
		}
		alerts = append(alerts, alert)
	}

	return alerts, rows.Err()
}
