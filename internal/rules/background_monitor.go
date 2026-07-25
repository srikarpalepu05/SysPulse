package rules

import "syspulse/internal/collector"

type Alert struct {
	ProcessName string
	PID         int32
	CPU         float64
	MemoryMB    float64
	Reason      string
	Severity    string
}

func Evaluate(snapshots []collector.ProcessSnapshot) []Alert {
	alerts := make([]Alert, 0)

	for _, snapshot := range snapshots {
		if snapshot.CPU < 15 && snapshot.MemoryMB < 200 {
			continue
		}

		severity := "low"
		reason := "Background process usage detected"

		switch {
		case snapshot.CPU >= 35 || snapshot.MemoryMB >= 900:
			severity = "high"
			reason = "High CPU or memory usage while running in the background"
		case snapshot.CPU >= 20 || snapshot.MemoryMB >= 350:
			severity = "medium"
			reason = "Notable background resource usage"
		default:
			severity = "low"
			reason = "Background process is using measurable resources"
		}

		alerts = append(alerts, Alert{
			ProcessName: snapshot.Name,
			PID:         snapshot.PID,
			CPU:         snapshot.CPU,
			MemoryMB:    snapshot.MemoryMB,
			Reason:      reason,
			Severity:    severity,
		})
	}

	return alerts
}
