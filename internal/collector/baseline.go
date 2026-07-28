package collector

import (
	"fmt"
	"time"
)

func BuildBaselineSummary(startupCount, processCount int) string {
	return fmt.Sprintf("baseline: %d startup items, %d active processes", startupCount, processCount)
}

func CaptureBaseline(startupEntries []StartupEntry, snapshots []ProcessSnapshot) map[string]interface{} {
	return map[string]interface{}{
		"generated_at":  time.Now(),
		"startup_count": len(startupEntries),
		"process_count": len(snapshots),
		"summary":       BuildBaselineSummary(len(startupEntries), len(snapshots)),
	}
}
