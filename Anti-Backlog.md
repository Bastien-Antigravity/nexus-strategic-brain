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
# 🚫 The Anti-Backlog

> "The historical record of conscious decisions NOT to implement a feature, pattern, or library."

## 🎯 Purpose
To prevent "Architectural Amnesia" where the same discarded ideas are re-debated every few months. If an idea is in the Anti-Backlog, it requires a significant change in ecosystem context to be reconsidered.

## 📁 Discarded Ideas

### 1. Direct Socket.io for Core Logging
- **Decision Date**: 2026-04-20
- **Reason**: Performance overhead and lack of memory safety in the node-js implementation for high-frequency logs.
- **Alternative**: Custom **SafeSocket** protocol over raw TCP.

### 2. Multi-Runtime FFI (Multiple Shared Libs)
- **Decision Date**: 2026-05-02
- **Reason**: Caused "Multiple Go Runtime" panics and memory isolation issues when Python/Rust loaded separate Go-based shared libraries.
- **Alternative**: The **Super-Bridge** (`universal-logger`).

### 3. Fragmented Knowledge Silos (PARA-only)
- **Decision Date**: 2026-05-16
- **Reason**: Caused "Reasoning Lag" and high token consumption due to excessive context-switching between 10+ repositories.
- **Alternative**: **Engineering-Brain Consolidation** (Merging Tech-Stack and Core-KMS).



---
*Reference: [[01-Strategic-Nexus/README]], [[Global-Architecture-Rules]]*
