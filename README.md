# SysPulse — Local Security, Hygiene, and System Intelligence Monitor

> [!IMPORTANT]
> SysPulse, NetScope, and PacketTrail have been moved from a previous GitHub account as part of a cleanup and reorganization effort.

SysPulse is a local system health and security project built to monitor what is happening on a machine, flag unusual behavior, and help keep the system clean, healthy, and trustworthy over time. It is designed to work like a practical local security and optimization tool, not just a dashboard for CPU and memory numbers.

## What this project is trying to do
The real goal is to help a machine understand its own behavior. SysPulse watches for:

- suspicious startup items
- drift away from a known-good baseline
- unusual process patterns
- weak system hygiene
- repeated risk signals that could point to a security or optimization issue

This is not just monitoring for performance. It is geared toward proactive local system awareness and security hygiene.

## Current status

### ✅ Core runtime is active in Go
The main system is currently built around Go and is responsible for:

- collecting live process snapshots
- scanning Windows startup entries
- storing baseline data in SQLite
- tracking drift and hygiene checks
- evaluating suspicious background activity
- saving historical alerts and reports

This is the active foundation of the project right now.

### ✅ Python intelligence has started
A Python-based intelligence layer is also being created to read the stored data and summarize system health in a more readable way. It is currently focused on:

- hygiene score summaries
- risk trend checks
- top process discovery
- startup item counts
- recent alert summaries
- action-oriented reporting

This intelligence layer is still in its early stage, but it is now an actual part of the project instead of just being an idea.

### ⚠️ UI is still in progress
The project is prioritizing the working logic first. A polished desktop experience is not the main goal right now. The important thing is to make the monitoring, scoring, and data collection reliable before spending time on visual polish.

## What is already working

- Process monitoring
- Startup inventory collection
- Baseline creation
- Drift detection against a saved baseline
- Hygiene scoring and risk summaries
- Recommendations based on detected issues
- Alert and report persistence in SQLite
- Python report generation from stored data

## What is next

### Phase 1: strengthen the core monitor
This is already the main working layer and should continue to be improved.

Planned next steps:
- improve console output readability
- make the health summary cleaner and more dashboard-like
- refine risk labels and severities
- improve baseline lifecycle management
- improve retention and deduplication of recorded data

### Phase 2: build the Python intelligence layer more seriously
This is the next proper phase of the project.

Planned next steps:
- summarize recent alert severity trends
- compare current behavior against historical trends
- detect startup risk patterns over time
- create a single “security pulse” summary for the terminal
- export reports in a format that can later be used for dashboards or automation

### Phase 3: stronger anomaly detection
Once the local monitor and intelligence layer are stable, the more advanced logic can be added.

Possible directions:
- detect abnormal CPU spikes that do not match normal behavior
- identify suspicious process patterns across time periods
- score startup items based on command, source, and risk pattern
- add more context-aware recommendations instead of generic warnings

### Phase 4: user-facing presentation
After the system logic is solid, the project can grow into a more visual layer.

Possible directions:
- a desktop dashboard
- a lightweight web dashboard
- a local admin panel
- notifications for risky system patterns

### Phase 5: later expansion
After the local product works well, the project can grow into a broader security system.

Possible directions:
- remote agents
- distributed monitoring across machines
- stronger centralized reports
- broader automation and security workflows

## Tech stack

- Go for the active monitor and runtime logic
- SQLite for local persistence and historical data
- Python for analytics and intelligence reporting
- local monitoring patterns built around process and startup behavior

## Project structure

```text
SysPulse/
├── cmd/app/                # main Go monitor loop
├── internal/collector/     # process and startup collectors
├── internal/analysis/      # hygiene, drift, anomaly, and risk logic
├── internal/rules/         # alerting and startup heuristics
├── internal/storage/       # database and persistence
├── python_intelligence/    # Python-based reporting and analytics
├── tests/                 # Python validation tests
├── go.mod                 # Go module definition
├── syspulse.db            # local SQLite database
├── run.py                 # legacy prototype reference
├── README.md              # project overview and roadmap
```



## Contact

📧 [srikarpalepu05@gmail.com](mailto:srikarpalepu05@gmail.com)


