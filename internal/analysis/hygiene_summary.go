package analysis

import "fmt"

func SecurityHygieneSummary(score int, status string) string {
	switch status {
	case "critical":
		return fmt.Sprintf("System hygiene is critical (score: %d/100)", score)
	case "warning":
		return fmt.Sprintf("System hygiene is warning (score: %d/100)", score)
	default:
		return fmt.Sprintf("System hygiene is healthy (score: %d/100)", score)
	}
}
