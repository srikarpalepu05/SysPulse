from agents.metrics_collector import collect_metrics
from storage.models import init_db
from storage.database import get_connection

def main():
    init_db()
    metrics = collect_metrics()

    conn = get_connection()
    cursor = conn.cursor()

    cursor.execute(
        "INSERT INTO system_metrics (timestamp, cpu, memory) VALUES (?, ?, ?)",
        (metrics["timestamp"], metrics["cpu"], metrics["memory"])
    )

    conn.commit()
    conn.close()

    print("Collected metrics:", metrics)

if __name__ == "__main__":
    main()
