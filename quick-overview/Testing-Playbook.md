--- 
title: Testing Playbook - Strategic Nexus
microservice: nexus-strategic-brain
type: human-doc
status: active
tags:
- '#zone/3-fleet'
- '#service/nexus-strategic-brain'
- '#state/active'
- '#type/human-doc'
- '#ai/ignore'
---
# 🧪 Testing Playbook: Strategic Nexus

This document provides a human-friendly summary of how we ensure the integrity of the **Strategic Nexus**. Unlike code-heavy microservices, "testing" here focuses on **reasoning consistency** and **link health**.

## 🛡️ How we verify Strategy
We use the **Sentinel (Role 09)** to audit the Oracle's output. The core goal is to ensure that every strategic claim ("Truth Bomb") is backed by evidence in the logs.

### 1. The "Truth Bomb" Check
Whenever a strategic audit is published, we verify:
- **Evidence**: Is there a link to the specific session or code that caused the concern?
- **Actionability**: Does it provide a clear "Next Step" or a "Global Law"?

### 2. The Anti-Backlog Guard
We check that new implementation plans do not accidentally "re-invent" something we already rejected in the **Anti-Backlog**. This prevents circular debates and wasted time.

### 3. Automated Health Checks
We run periodic scans to ensure:
- **Link Health**: No broken Obsidian links between strategy docs.
- **Metadata**: Every file is correctly tagged for the Dataview dashboards.

## ⚙️ Key Rituals for Humans
If you are manually auditing the brain, you can run:
- **Link Validator**: `python3 08-Base-Scripts/archive/vault-sentinel.py --path 01-Strategic-Nexus`
- **Health Audit**: `python3 08-Base-Scripts/Brain-Health-Audit.py`

---
> **Authority**: For the full technical specification of these rituals, see: **[[01-Level-Governance#🧪 5. Quality & Verification (The Sentinel Audit)|Quality & Verification (Full Spec)]]**
