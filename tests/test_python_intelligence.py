import sqlite3
from pathlib import Path

from python_intelligence.insights import compute_risk_summary


def test_compute_risk_summary_reads_history_and_reports_trend(tmp_path):
    db_path = tmp_path / "syspulse.db"
    conn = sqlite3.connect(db_path)
    conn.execute(
        """
        CREATE TABLE hygiene_reports (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            timestamp TEXT NOT NULL,
            score INTEGER NOT NULL,
            status TEXT NOT NULL,
            reason TEXT NOT NULL,
            risk_breakdown TEXT NOT NULL,
            recommendation TEXT NOT NULL
        )
        """
    )
    conn.execute(
        """
        CREATE TABLE process_snapshots (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            timestamp TEXT NOT NULL,
            process_name TEXT NOT NULL,
            pid INTEGER NOT NULL,
            cpu_percent REAL NOT NULL,
            memory_mb REAL NOT NULL
        )
        """
    )
    conn.execute(
        """
        CREATE TABLE startup_items (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            timestamp TEXT NOT NULL,
            name TEXT NOT NULL,
            command TEXT NOT NULL,
            source TEXT NOT NULL,
            location TEXT NOT NULL
        )
        """
    )
    conn.executemany(
        "INSERT INTO hygiene_reports (timestamp, score, status, reason, risk_breakdown, recommendation) VALUES (?, ?, ?, ?, ?, ?)",
        [
            ("2026-08-09T00:00:00Z", 60, "warning", "baseline drift", "medium", "investigate"),
            ("2026-08-09T00:00:10Z", 78, "good", "stable", "low", "monitor"),
            ("2026-08-09T00:00:20Z", 72, "warning", "minor drift", "medium", "review"),
        ],
    )
    conn.executemany(
        "INSERT INTO process_snapshots (timestamp, process_name, pid, cpu_percent, memory_mb) VALUES (?, ?, ?, ?, ?)",
        [
            ("2026-08-09T00:00:00Z", "chrome", 101, 30.0, 400.0),
            ("2026-08-09T00:00:05Z", "chrome", 101, 50.0, 480.0),
            ("2026-08-09T00:00:10Z", "python", 202, 80.0, 700.0),
        ],
    )
    conn.executemany(
        "INSERT INTO startup_items (timestamp, name, command, source, location) VALUES (?, ?, ?, ?, ?)",
        [
            ("2026-08-09T00:00:00Z", "PowerShell", "powershell.exe", "registry", "HKCU\\Startup"),
            ("2026-08-09T00:00:00Z", "Discord", "Discord.exe", "registry", "HKCU\\Startup"),
        ],
    )
    conn.commit()
    conn.close()

    summary = compute_risk_summary(str(db_path))

    assert summary["latest_score"] == 72
    assert summary["average_score"] == 70.0
    assert summary["trend"] in {"improving", "stable", "declining"}
    assert "chrome" in summary["top_processes"]
    assert summary["startup_items_count"] == 2
    assert summary["risk_level"] in {"low", "moderate", "high"}


def test_compute_risk_summary_handles_missing_tables(tmp_path):
    db_path = tmp_path / "missing.db"
    conn = sqlite3.connect(db_path)
    conn.close()

    summary = compute_risk_summary(str(db_path))

    assert summary["latest_score"] == 0
    assert summary["average_score"] == 0.0
    assert summary["trend"] == "stable"
    assert summary["top_processes"] == []
    assert summary["startup_items_count"] == 0
    assert summary["risk_level"] == "low"
