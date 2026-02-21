# SysPulse — Intelligent System Health, Security & Performance Monitoring Platform (In Progress)

> [!IMPORTANT]  
> **NOTE:** SysPulse, NetScope, and PacketTrail have been moved from my previous GitHub account (`srikar-palepu`) as part of a cleanup and reorganization effort.

SysPulse is an evolving system health, security, and performance monitoring platform designed to collect, store, and analyze system metrics across multiple machines. The project is currently under active development.

## 🚀 Current Capabilities
- **✔ Local System Metrics Collection**
  - CPU usage, memory usage, disk activity, network I/O
  - Built using [psutil](https://psutil.readthedocs.io)
  - Modular collectors for easy expansion
- **✔ Persistent Local Storage**
  - SQLite backend
  - Efficient schema for historical querying
  - Supports emerging analytics workflows
- **✔ Modular Architecture**
  - Clean folder structure
  - Easy to integrate new collectors, alerting modules, or external APIs

## 🛠️ Planned Features (Under Development)
- **🔹 Distributed Agent Support**
  - Lightweight agents running remotely and sending metric streams to a central server.
- **🔹 Threshold-Based Alerting**
  - Custom triggers for: CPU spikes, memory leaks, and unusual network activity.
- **🔹 Historical Analytics Dashboard**
  - Long-term trends and insights for performance and reliability.
- **🔹 AI-Assisted Anomaly Detection (Upcoming)**
  - Using [scikit-learn](https://scikit-learn.org) models for pattern recognition and early failure warnings.
  - *Note: AI features are in prototyping stage and not yet fully integrated.*
- **🔹 Cloud-Ready Architecture**
  - Future upgrades aimed at centralized ingestion APIs, scalable storage, and real-time web dashboards.

## 🧪 Tech Stack
*   **Core:** [Python](https://www.python.org) & [SQLite](https://www.sqlite.org)
*   **Monitoring:** [psutil](https://github.com)
*   **Data Science (Planned):** [Pandas](https://pandas.pydata.org) & [scikit-learn](https://scikit-learn.org)
*   **Infrastructure (Planned):** [Flask](https://flask.palletsprojects.com) & [Docker](https://www.docker.com)

## 📌 Project Status
> SysPulse is a **work in progress**, actively expanding into system health analytics, distributed monitoring, and AI-based anomaly detection.

## 📫 Contact
For questions or collaboration:  
📧 **[srikarpalepu05@gmail.com](mailto:srikarpalepu05@gmail.com)**


## 📁 Project Structure
```text
sysPulse/
│── agents/        # Remote collectors (in progress)
│── storage/       # SQLite schema + migrations
│── core/          # Ingestion + processing
│── analysis/      # ML + anomaly detection (future)
│── cli/           # Command-line utilities (future)


