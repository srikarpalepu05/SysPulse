package main

import (
	"fmt"
	"log"
	"time"

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

	fmt.Println("SysPulse Background Monitor started")
	fmt.Println("Monitoring process activity and startup items every 10 seconds. Press Ctrl+C to stop.\n")

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
			for _, entry := range startupEntries {
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
			fmt.Printf("Startup entries captured: %d\n", len(startupEntries))
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
