package analysis

import "fmt"

type DriftRecommendation struct {
	Type     string
	Severity string
	Action   string
	Reason   string
}

func DriftRecommendations(findings []DriftFinding) []DriftRecommendation {
	recommendations := make([]DriftRecommendation, 0, len(findings))

	for _, finding := range findings {
		action := "Continue monitoring and review the change during the next maintenance window."
		reason := "The current state differs from the saved baseline and should be inspected."

		switch finding.Type {
		case "startup":
			action = "Review new or unexpected startup entries and remove any untrusted programs."
			reason = "Startup items changed from the baseline and may indicate a new persistence mechanism."
		case "process":
			action = "Inspect newly running background processes and verify they are expected and legitimate."
			reason = "The active process count changed from the baseline and may indicate a resource or security drift."
		}

		switch finding.Severity {
		case "high":
			action = "Immediately investigate and remove unexpected startup entries or suspicious processes."
		case "medium":
			action = "Review the change soon and verify whether the drift is intentional or risky."
		}

		recommendations = append(recommendations, DriftRecommendation{
			Type:     finding.Type,
			Severity: finding.Severity,
			Action:   action,
			Reason:   reason,
		})
	}

	return recommendations
}

func DriftSummary(findings []DriftFinding) string {
	if len(findings) == 0 {
		return "System is within baseline"
	}

	highCount := 0
	for _, finding := range findings {
		if finding.Severity == "high" {
			highCount++
		}
	}

	if highCount > 0 {
		return fmt.Sprintf("System drift detected: %d high-severity findings", highCount)
	}

	return fmt.Sprintf("System drift detected: %d findings requiring review", len(findings))
}
