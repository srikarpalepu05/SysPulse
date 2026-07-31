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
            "recent_alerts": [],
            "needs_attention": "No data yet. Start the monitor to build history.",
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

        alert_table_exists = conn.execute(
            "SELECT 1 FROM sqlite_master WHERE type='table' AND name='alerts'"
        ).fetchone()
        if alert_table_exists:
            alert_rows = conn.execute(
                """
                SELECT process_name, severity, reason
                FROM alerts
                ORDER BY id DESC
                LIMIT 5
                """
            ).fetchall()
        else:
            alert_rows = []
    except sqlite3.Error:
        conn.close()
        return {
            "latest_score": 0,
            "average_score": 0.0,
            "trend": "stable",
            "top_processes": [],
            "startup_items_count": 0,
            "risk_level": "low",
            "recent_alerts": [],
            "needs_attention": "No data yet. Start the monitor to build history.",
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

    recent_alerts = [
        {
            "process_name": row["process_name"],
            "severity": row["severity"],
            "reason": row["reason"],
        }
        for row in alert_rows
    ]

    if risk_level == "high":
        needs_attention = "Immediate review recommended: check high-risk processes and startup items."
    elif risk_level == "moderate":
        needs_attention = "Monitor closely and review suspicious processes before they escalate."
    else:
        needs_attention = "System looks stable; keep monitoring for unusual changes."

    return {
        "latest_score": latest_score,
        "average_score": average_score,
        "trend": trend,
        "top_processes": top_processes,
        "startup_items_count": startup_rows["total"] if startup_rows else 0,
        "risk_level": risk_level,
        "recent_alerts": recent_alerts,
        "needs_attention": needs_attention,
    }


def build_risk_report(db_path: str) -> str:
    """Render a text summary suitable for a terminal dashboard."""
    summary = compute_risk_summary(db_path)
    top_processes = ", ".join(summary["top_processes"]) if summary["top_processes"] else "none"
    alerts = summary["recent_alerts"]

    alert_lines = [
        f"- {alert['process_name']} [{alert['severity']}]: {alert['reason']}"
        for alert in alerts
    ] or ["- none"]

    return (
        "SysPulse Python Intelligence Report\n"
        "=================================\n"
        f"Latest hygiene score: {summary['latest_score']}\n"
        f"Average hygiene score: {summary['average_score']}\n"
        f"Trend: {summary['trend']}\n"
        f"Risk level: {summary['risk_level']}\n"
        f"Startup items: {summary['startup_items_count']}\n"
        f"Top processes: {top_processes}\n\n"
        "Recent alerts:\n"
        + "\n".join(alert_lines)
        + "\n\n"
        + f"Needs attention: {summary['needs_attention']}\n"
    )


def main() -> None:
    parser = argparse.ArgumentParser(description="Generate a SysPulse intelligence summary from the local database.")
    parser.add_argument("--db-path", default="syspulse.db", help="Path to the SysPulse SQLite database.")
    args = parser.parse_args()

    print(build_risk_report(args.db_path))


if __name__ == "__main__":
    main()
