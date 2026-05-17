---
microservice: nexus-strategic-brain
type: architecture
status: active
tags:
- \'#zone/3-fleet\'
- '#ai/ignore'
- '#service/nexus-strategic-brain'
---

# 📚 Strategic Nexus: General & Misc

This document covers the "Librarian" and meta-data standards that keep the Strategic Nexus navigable and machine-readable for the AI Squad.

## 🏷️ Tagging & Metadata Standards

Every document in the Nexus must contain a frontmatter block with at least the following:
- `microservice`: `nexus-strategic-brain`
- `type`: One of `meta-audit`, `architecture`, `strategy`, or `documentation`.
- `status`: `active`, `superseded`, or `archived`.

### Specialized Tags
- `#ai/ignore`: Applied to `quick-overview/` to prevent execution agents (Developers) from being overwhelmed by high-level strategic context during code-centric tasks.
- `#type/meta-audit`: Specifically for `STRAT-XXX` files.
- `#state/truth-bomb`: For specific observations that require immediate human attention.

## 📂 File Naming Conventions

- **Strategic Audits**: `STRAT-[NNN]-[Descriptive-Name].md` (e.g., `STRAT-003-The-Infrastructure-Gordian-Knot.md`).
- **Philosophy**: `PHIL-[NNN]-[Topic].md`.
- **Patterns**: Functional names (e.g., `Strategic-Patterns.md`).

## 📊 Knowledge Synthesis (Dataview)

The Nexus leverages the **Obsidian Dataview** plugin to create live dashboards. 
- **The Strategy Index**: Automatically tracks all `STRAT` documents, their status, and dates.
- **Dependency Maps**: Visualizes how "Rejected" ideas in the Anti-Backlog relate to current "Active" architectural rules.

## 🧠 Philosophy: "Slow is Smooth"

The Nexus operates on a different time-scale than the rest of the fleet.
- **Rule**: If a strategic document is written in haste, it is likely just another execution log.
- **Guideline**: Strategic synthesis should only happen after a minimum of 3 development sessions to ensure enough "Log Gravity" has accumulated to form a clear pattern.

## 🔗 Repository Links

- **Main Repository**: [01-Strategic-Nexus](https://github.com/Bastien-Antigravity/01-Strategic-Nexus)
- **Primary Role**: [[07-Core-KMS/Role-Prompts/00-Oracle/Prompt-Chronos-Oracle|Chronos-Oracle]]
- **Parent Hub**: [[Ecosystem-Map-MOC]]

