package analysis

import "testing"

func TestFixRecommendationForRiskyStartup(t *testing.T) {
	breakdown := RiskBreakdown{SuspiciousStartup: 36, ProcessRisk: 0, StartupRisk: 0}
	rec := FixRecommendation(breakdown, 60)
	if rec == "" {
		t.Fatal("expected recommendation text")
	}
	if rec != "Review and remove unexpected startup entries, especially scripts or command shells." {
		t.Fatalf("unexpected recommendation: %q", rec)
	}
}

func TestFixRecommendationSummaryIncludesPrefix(t *testing.T) {
	summary := FixRecommendationSummary(RiskBreakdown{}, 90)
	if summary != "Recommended fix: No urgent action needed; continue monitoring and keep the baseline stable." {
		t.Fatalf("unexpected summary: %q", summary)
	}
}
