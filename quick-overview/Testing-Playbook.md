---
microservice: nexus-strategic-brain
type: architecture
status: active
tags:
- '#ai/ignore'
- '#service/nexus-strategic-brain'
---

# 🧪 Strategic Nexus: Testing Playbook

Testing in the Strategic Nexus is not about "Unit Tests" but about **Semantic Integrity** and **Reasoning Consistency**. We test the Oracle to ensure its "Prophecies" are grounded in reality.

## 🛡️ The Sentinel Audit

The **Sentinel** (Role 09) is the primary auditor for this repository. It runs automated and manual checks based on the following rules:

### 1. Structural Integrity (Hard Checks)
- **Link Integrity**: Every `STRAT-XXX` document must have working links to the logs or code repositories it critiques.
- **Tag Compliance**: Every file must have valid frontmatter (no `null` tags, correct `microservice` ID).
- **File Locality**: Strategic documents must remain in `01-Strategic-Nexus` and not leak into execution-focused folders.

### 2. Semantic Consistency (Soft Checks)
- **The "Truth Bomb" Validaton**: If a `STRAT-XXX` audit makes a claim (e.g., "SafeSocket has an OOM risk"), the Sentinel verifies that there is supporting evidence in the `sandbox-testing` logs or session history.
- **Anti-Backlog Conflict**: Ensures that a new feature proposal does not conflict with a previously settled rejection in the Anti-Backlog without a formal "Strategic Pivot" document.

## 🧪 Quality Assurance Workflows

### The "Loop-Back" Test
1.  **Input**: Select a recommendation from a previous `STRAT` document (e.g., "Implement stateless bootstrapping").
2.  **Verification**: Search the fleet (using `fleet-manager.py`) to see if the recommendation was actually implemented.
3.  **Result**: If not implemented, the Nexus must generate a "Follow-up Audit" to investigate the execution bottleneck.

### Link Integrity Validation
Run the centralized link validator to ensure that the cross-repo knowledge graph is not broken:
```bash
python3 20-Scripts/verify_links_script.py --path 01-Strategic-Nexus
```

## 📈 Success Metrics

- **Drift Detection**: The time between a "Blind Spot" occurring in code and it being captured in a `STRAT` audit.
- **Rejection Rate**: The number of recurring technical debates silenced by the Anti-Backlog.
- **Oracle Accuracy**: The ratio of strategic predictions that correctly identified infrastructure bottlenecks before they caused a "Fleet Deadlock."

## 🤖 Sentinel Prompt Snippet
> "As a Sentinel, your duty is to audit the Oracle. Ensure that every 'Why' has a corresponding 'Where' in the logs, and every 'What' has a corresponding 'When' in the project timeline."

