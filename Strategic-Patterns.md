---

microservice: nexus-strategic-brain
type: strategy
status: active
tags:
- \'#service/nexus-strategic-brain\'
- '#type/strategy'
- '#state/active'
- '#tier/strategy'

---
# 🧩 Strategic Patterns

> "Recurring architectural and operational truths that govern the Bastien-Antigravity ecosystem."

## 🏗️ Architectural Patterns
- **The Facade Law**: Core logic must always be wrapped in a language-agnostic facade (FFI/gRPC) to ensure polyglot compatibility.
- **Static Loading Law**: Binaries must be statically linked to avoid shared object hell in Docker environments.
- **The Super-Bridge**: Consolidated shared infrastructure (Logger + Config) into a single FFI boundary to prevent runtime panics.

## 🚀 Operational Patterns
- **Spec-First Enforcement**: Zero coding starts without an approved BDD spec.
- **Fail-Fast Security**: Reject unauthenticated or oversized payloads at the socket layer, not the application layer.
- **Atomic Fleet Actions**: If a multi-repo update fails on one repo, the entire action is halted for manual review.
- **The Sovereignty Ritual**: Mandatory Mission Sign-off and Git-aware stateless auditing before any fleet-wide merge.
- **Multi-AI Agnostic CLI**: All orchestration scripts must run on standard `python3` and avoid agent-specific syntax to ensure cross-platform/cross-AI compatibility.


---
*Reference: [[01-Strategic-Nexus/README]], [[STRAT-001-The-Dormant-Pipeline]]*
