from __future__ import annotations

import argparse
import sqlite3
from collections import Counter
from statistics import mean
from typing import Any


def compute_risk_summary(db_path: str) -> dict[str, Any]:
    """Read recent SysPulse data and summarize current security risk."""
    try:
        conn = sqlite3.connect(db_path)
    except sqlite3.Error:
        return {
            "latest_score": 0,
            "average_score": 0.0,
            "trend": "stable",
            "top_processes": [],
            "startup_items_count": 0,
            "risk_level": "low",
        }

    conn.row_factory = sqlite3.Row

    try:
        latest_hygiene = conn.execute(
            """
            SELECT score, status, reason
            FROM hygiene_reports
            ORDER BY id DESC
            LIMIT 1
            """
        ).fetchone()

        score_history = conn.execute(
            """
            SELECT score
            FROM hygiene_reports
            ORDER BY id ASC
            """
        ).fetchall()

        process_rows = conn.execute(
            """
            SELECT process_name, cpu_percent, memory_mb
            FROM process_snapshots
            ORDER BY id DESC
            """
        ).fetchall()

        startup_rows = conn.execute(
            "SELECT COUNT(*) AS total FROM startup_items"
        ).fetchone()
    except sqlite3.Error:
        conn.close()
        return {
            "latest_score": 0,
            "average_score": 0.0,
            "trend": "stable",
            "top_processes": [],
            "startup_items_count": 0,
            "risk_level": "low",
        }

    latest_score = latest_hygiene["score"] if latest_hygiene else 0
    history_scores = [row["score"] for row in score_history]
    average_score = round(mean(history_scores), 2) if history_scores else 0.0

    if len(history_scores) >= 2:
        delta = history_scores[-1] - history_scores[0]
        if delta > 5:
            trend = "improving"
        elif delta < -5:
            trend = "declining"
        else:
            trend = "stable"
    else:
        trend = "stable"

    process_usage = Counter()
    for row in process_rows:
        process_usage[row["process_name"]] += 1

    conn.close()

    top_processes = [name for name, _ in process_usage.most_common(5)]

    if latest_score >= 80:
        risk_level = "low"
    elif latest_score >= 60:
        risk_level = "moderate"
    else:
        risk_level = "high"

    return {
        "latest_score": latest_score,
        "average_score": average_score,
        "trend": trend,
        "top_processes": top_processes,
        "startup_items_count": startup_rows["total"] if startup_rows else 0,
        "risk_level": risk_level,
    }


def build_risk_report(db_path: str) -> str:
    """Render a text summary suitable for a terminal dashboard."""
    summary = compute_risk_summary(db_path)
    top_processes = ", ".join(summary["top_processes"]) if summary["top_processes"] else "none"

    return (
        "SysPulse Python Intelligence Report\n"
        "=================================\n"
        f"Latest hygiene score: {summary['latest_score']}\n"
        f"Average hygiene score: {summary['average_score']}\n"
        f"Trend: {summary['trend']}\n"
        f"Risk level: {summary['risk_level']}\n"
        f"Startup items: {summary['startup_items_count']}\n"
        f"Top processes: {top_processes}\n"
    )


def main() -> None:
    parser = argparse.ArgumentParser(description="Generate a SysPulse intelligence summary from the local database.")
    parser.add_argument("--db-path", default="syspulse.db", help="Path to the SysPulse SQLite database.")
    args = parser.parse_args()

    print(build_risk_report(args.db_path))


if __name__ == "__main__":
    main()
