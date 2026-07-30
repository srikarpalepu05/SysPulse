package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSaveAndGetLatestHygieneReport(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "hygiene.db")

	db, err := OpenDB(path)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := EnsureHygieneTable(db); err != nil {
		t.Fatalf("ensure hygiene table: %v", err)
	}

	report := HygieneReport{
		Timestamp:      time.Now(),
		Score:          81,
		Status:         "healthy",
		Reason:         "System is stable",
		RiskBreakdown:  "startup drift: 12 | process drift: 6",
		Recommendation: "Review startup entries",
	}

	if err := SaveHygieneReport(db, report); err != nil {
		t.Fatalf("save hygiene report: %v", err)
	}

	latest, err := GetLatestHygieneReport(db)
	if err != nil {
		t.Fatalf("get latest hygiene report: %v", err)
	}
	if latest.Score != report.Score {
		t.Fatalf("expected score %d, got %d", report.Score, latest.Score)
	}
	if latest.Status != report.Status {
		t.Fatalf("expected status %s, got %s", report.Status, latest.Status)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("database file not created: %v", err)
	}
}
