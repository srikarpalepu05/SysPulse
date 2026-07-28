package rules

import "fmt"

type Recommendation struct {
	Name     string
	Source   string
	Action   string
	Reason   string
	Severity string
}

func StartupRecommendations(risks []StartupRisk) []Recommendation {
	recommendations := make([]Recommendation, 0, len(risks))

	for _, risk := range risks {
		action := "Monitor this startup item and review it during the next maintenance window."
		reason := "This startup item has low to moderate risk and may be safe to keep running."

		switch risk.Severity {
		case "high":
			action = fmt.Sprintf("Disable this startup item or review %s before it launches automatically.", risk.Name)
			reason = "High-risk startup pattern detected. This entry may be consuming resources or increasing exposure."
		case "medium":
			action = fmt.Sprintf("Review %s and consider delaying or removing it from startup if it is not required.", risk.Name)
			reason = "This item matches a moderate-risk pattern and is worth checking."
		default:
			action = fmt.Sprintf("Keep monitoring %s, but it is currently a low-priority startup item.", risk.Name)
			reason = "This item is only mildly suspicious and should be watched rather than removed immediately."
		}

		recommendations = append(recommendations, Recommendation{
			Name:     risk.Name,
			Source:   risk.Source,
			Action:   action,
			Reason:   reason,
			Severity: risk.Severity,
		})
	}

	return recommendations
}
