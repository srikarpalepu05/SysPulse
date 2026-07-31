from __future__ import annotations

import sqlite3
from collections import Counter
from statistics import mean
from typing import Any


def compute_risk_summary(db_path: str) -> dict[str, Any]:
    """Read recent SysPulse data and summarize current security risk."""
    conn = sqlite3.connect(db_path)
    conn.row_factory = sqlite3.Row

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

    process_rows = conn.execute(
        """
        SELECT process_name, cpu_percent, memory_mb
        FROM process_snapshots
        ORDER BY id DESC
        """
    ).fetchall()

    process_usage = Counter()
    for row in process_rows:
        process_usage[row["process_name"]] += 1

    startup_rows = conn.execute(
        "SELECT COUNT(*) AS total FROM startup_items"
    ).fetchone()

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
