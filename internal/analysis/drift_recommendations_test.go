package analysis

import "testing"

func TestDriftRecommendationsForStartupChange(t *testing.T) {
	findings := []DriftFinding{
		{Type: "startup", Severity: "medium", Message: "startup count changed from 3 to 5", Baseline: 3, Current: 5, Count: 2},
	}

	recs := DriftRecommendations(findings)
	if len(recs) == 0 {
		t.Fatal("expected recommendations for startup drift")
	}

	if recs[0].Severity != "medium" {
		t.Fatalf("expected medium severity recommendation, got %s", recs[0].Severity)
	}
}

func TestDriftSummaryForNoFindings(t *testing.T) {
	text := DriftSummary(nil)
	if text != "System is within baseline" {
		t.Fatalf("expected baseline summary, got %q", text)
	}
}
