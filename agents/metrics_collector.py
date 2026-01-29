import psutil
from datetime import datetime

def collect_metrics():
    return {
        "timestamp": datetime.utcnow().isoformat(),
        "cpu": psutil.cpu_percent(),
        "memory": psutil.virtual_memory().percent
    }
