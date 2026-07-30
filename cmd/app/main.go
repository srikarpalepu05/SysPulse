package main

import (
	"fmt"
	"log"
	"time"

	"syspulse/internal/analysis"
	"syspulse/internal/collector"
	"syspulse/internal/rules"
	"syspulse/internal/storage"
)

func main() {
	db, err := storage.OpenDB("syspulse.db")
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	if err := storage.EnsureStartupTable(db); err != nil {
		log.Fatalf("ensure startup table: %v", err)
	}
	if err := storage.EnsureBaselineTable(db); err != nil {
		log.Fatalf("ensure baseline table: %v", err)
	}
	fmt.Println("SysPulse Background Monitor started")
	fmt.Println("Monitoring process activity every 10 seconds. Press Ctrl+C to stop.")

	for {
		snapshots, err := collector.Collect()
		if err != nil {
			fmt.Println("Failed to collect process data:", err)
			time.Sleep(10 * time.Second)
			continue
		}

		for _, snapshot := range snapshots {
			if err := storage.SaveSnapshot(db, snapshot); err != nil {
				log.Printf("save snapshot: %v", err)
			}
		}

		alerts := rules.Evaluate(snapshots)
		for _, alert := range alerts {
			if err := storage.SaveAlert(db, alert); err != nil {
				log.Printf("save alert: %v", err)
			}
		}

		startupEntries, err := collector.CollectStartupEntries()
		if err != nil {
			fmt.Println("Failed to collect startup entries:", err)
		} else {
			seenStartup := make(map[string]struct{}, len(startupEntries))
			for _, entry := range startupEntries {
				key := fmt.Sprintf("%s|%s|%s", entry.Source, entry.Name, entry.Command)
				if _, exists := seenStartup[key]; exists {
					continue
				}
				seenStartup[key] = struct{}{}

				if err := storage.SaveStartupEntry(db, storage.StartupRecord{
					Name:      entry.Name,
					Command:   entry.Command,
					Source:    entry.Source,
					Location:  entry.Location,
					Timestamp: entry.Timestamp,
				}); err != nil {
					log.Printf("save startup entry: %v", err)
				}
			}

			latestBaseline, err := storage.GetLatestBaseline(db)
			if err == nil {
				findings := analysis.DetectDrift(latestBaseline, startupEntries, snapshots)
				hygiene := analysis.ComputeHygieneScore(latestBaseline, startupEntries, snapshots)
				breakdown := analysis.ComputeRiskBreakdown(latestBaseline, startupEntries, snapshots)
				fmt.Println(analysis.SecurityHygieneSummary(hygiene.Score, hygiene.Status))
				fmt.Printf("Status: %s | %s\n", hygiene.Status, hygiene.Reason)
				fmt.Printf("Risk breakdown: %s\n", analysis.RiskBreakdownSummary(breakdown))
				fmt.Println(analysis.FixRecommendationSummary(breakdown, hygiene.Score))
				if len(findings) > 0 {
					fmt.Println("Baseline drift detected:")
					for _, finding := range findings {
						fmt.Printf("[%s] %s | baseline=%d current=%d\n", finding.Severity, finding.Message, finding.Baseline, finding.Current)
					}
					recommendations := analysis.DriftRecommendations(findings)
					for _, recommendation := range recommendations {
						fmt.Printf("[%s] %s — %s\n", recommendation.Severity, recommendation.Type, recommendation.Action)
					}
					fmt.Println(analysis.DriftSummary(findings))
				} else {
					fmt.Println(analysis.DriftSummary(nil))
				}
			}

			baseline := collector.CaptureBaseline(startupEntries, snapshots)
			baselineSnapshot := storage.BaselineSnapshot{
				GeneratedAt:  baseline["generated_at"].(time.Time),
				StartupCount: baseline["startup_count"].(int),
				ProcessCount: baseline["process_count"].(int),
				Summary:      baseline["summary"].(string),
			}
			if err := storage.SaveBaseline(db, baselineSnapshot); err != nil {
				log.Printf("save baseline: %v", err)
			}

			startupRisks := rules.ScoreStartupEntries(startupEntries)
			recommendations := rules.StartupRecommendations(startupRisks)
			fmt.Printf("Startup entries captured: %d | Risk-scored: %d | Recommendations: %d\n", len(startupEntries), len(startupRisks), len(recommendations))
			if len(startupRisks) == 0 {
				fmt.Println("No startup items matched the risk heuristics.")
			} else {
				for _, risk := range startupRisks {
					fmt.Printf("[%s] %s — score %d — %s — %s\n",
						risk.Severity,
						risk.Name,
						risk.Score,
						risk.Source,
						risk.Reason,
					)
				}
				for _, rec := range recommendations {
					fmt.Printf("[%s] %s — Action: %s\n",
						rec.Severity,
						rec.Name,
						rec.Action,
					)
				}
			}
		}

		fmt.Printf("Processes scanned: %d | Alerts: %d\n", len(snapshots), len(alerts))
		if len(alerts) == 0 {
			fmt.Println("No suspicious background activity detected.")
		} else {
			for _, alert := range alerts {
				fmt.Printf("[%s] %s (PID %d) — %.1f%% CPU — %.1f MB — %s\n",
					alert.Severity,
					alert.ProcessName,
					alert.PID,
					alert.CPU,
					alert.MemoryMB,
					alert.Reason,
				)
			}
		}

		fmt.Println("--------------------------------------------------")
		time.Sleep(10 * time.Second)
	}
}
