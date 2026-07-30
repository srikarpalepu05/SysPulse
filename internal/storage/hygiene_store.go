package storage

import (
	"database/sql"
	"fmt"
	"time"
)

type HygieneReport struct {
	Timestamp      time.Time
	Score          int
	Status         string
	Reason         string
	RiskBreakdown  string
	Recommendation string
}

func EnsureHygieneTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS hygiene_reports (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp TEXT NOT NULL,
			score INTEGER NOT NULL,
			status TEXT NOT NULL,
			reason TEXT NOT NULL,
			risk_breakdown TEXT NOT NULL,
			recommendation TEXT NOT NULL
		);`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("create hygiene_reports table: %w", err)
	}

	return nil
}

func SaveHygieneReport(db *sql.DB, report HygieneReport) error {
	query := `
		INSERT INTO hygiene_reports (timestamp, score, status, reason, risk_breakdown, recommendation)
		VALUES (?, ?, ?, ?, ?, ?)`

	_, err := db.Exec(query,
		report.Timestamp.Format(time.RFC3339),
		report.Score,
		report.Status,
		report.Reason,
		report.RiskBreakdown,
		report.Recommendation,
	)
	if err != nil {
		return fmt.Errorf("insert hygiene report: %w", err)
	}

	return nil
}

func GetLatestHygieneReport(db *sql.DB) (HygieneReport, error) {
	var report HygieneReport
	var ts string

	query := `
		SELECT timestamp, score, status, reason, risk_breakdown, recommendation
		FROM hygiene_reports
		ORDER BY id DESC
		LIMIT 1`

	row := db.QueryRow(query)
	if err := row.Scan(&ts, &report.Score, &report.Status, &report.Reason, &report.RiskBreakdown, &report.Recommendation); err != nil {
		return HygieneReport{}, err
	}

	parsed, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return HygieneReport{}, fmt.Errorf("parse hygiene report timestamp: %w", err)
	}
	report.Timestamp = parsed
	return report, nil
}
