---
microservice: nexus-strategic-brain
type: repository
status: active
tags:
- '#service/nexus-strategic-brain'
- '#state/active'
- '#tier/strategy'
- '#type/repository'
- '#zone/3-fleet'
---
# 🌌 Strategic Oracle Brain (01-Strategic-Nexus)

> "The Strategic Memory and Predictive Oracle of the Bastien-Antigravity Ecosystem."

## 🎯 Purpose
This repository acts as the **Meta-Intelligence Layer** for the project. It uses **Log-Driven Synthesis** to automatically analyze recent project history and session logs, focusing on **"Why"** and **"What Next"** to prevent reasoning drift and infrastructure procrastination.

## 📚 Ecosystem Documentation
This repository adheres to the fleet's documentation standards. Please refer to the following root files for core context:

- **[[quick-overview/Architecture-Overview|Architecture Overview]]**: The structural philosophy and 3-Tier Strategic Stack.
- **[[00-Level-Governance#🧪 5. Quality & Verification (The Sentinel Audit)|Quality & Verification]]**: The Sentinel Audit rules and validation workflows.
- **[[AI-Project-DNA]]**: High-level intent and AI role constraints.
- **[[AI-Session-State]]**: The current state of AI tasks in this repo.
- **[[TODO]]**: Pending tasks and technical debt.

*(Note: Human-friendly summaries are available in the `quick-overview/` directory).*

## 🧠 Strategic Memory (What we store)
To prevent "Architectural Amnesia," this brain tracks four specific types of meta-data:
1.  **STRAT-XXX (Strategic Audits)**: High-level reports on project trajectory and "Truth Bombs" regarding the current state.
2.  **Strategic Patterns**: Identification of "What always works" vs "What always breaks" in this specific ecosystem.
3.  **The Anti-Backlog**: A record of conscious decisions **NOT** to implement a feature or pattern, preventing recurring debates on settled topics.
4.  **Blind-Spot Logs**: Observations of risks that were invisible to the execution squad (Developers/Architects).

## 🛠️ Management CLI
The Strategic Nexus Go microservice automates role prompt scaffolding:
```bash
# Run the compiled binary to generate a new squad role prompt
./bin/strategic-nexus create-persona --name security-specialist --level 03 --specialty "Network Sec"
```


## 🔗 Connection
Linked to the **Lead Developer** loop via the `07-Core-KMS/Role-Prompts/00-Oracle/`.
