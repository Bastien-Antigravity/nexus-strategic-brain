#!/usr/bin/env python
# coding:utf-8

"""
🏗️ Persona Manager (Level 01 Strategic Nexus utility)
Handles the creation, template scaffolding, and registration of new AI Role Personas.
"""

import os
import sys
import argparse
import subprocess
from pathlib import Path

# --- Virtual Environment Re-Execution ---
_venv_dir = os.path.dirname(os.path.abspath(__file__))
while _venv_dir and _venv_dir != '/' and not os.path.exists(os.path.join(_venv_dir, ".venv")):
    _parent = os.path.dirname(_venv_dir)
    if _parent == _venv_dir:
        break
    _venv_dir = _parent
_venv_python = os.path.join(_venv_dir, ".venv", "bin", "python3") if os.name != "nt" else os.path.join(_venv_dir, ".venv", "Scripts", "python.exe")
if os.path.exists(_venv_python) and not os.path.samefile(sys.executable, _venv_python):
    try:
        os.execl(_venv_python, _venv_python, *sys.argv)
    except OSError:
        pass

def get_next_role_id(kms_prompts_dir: Path) -> int:
    """Scans 07-Core-KMS/Role-Prompts for folders prefixed with numbers to find the next ID."""
    max_id = -1
    if not kms_prompts_dir.exists():
        return 13  # Default fallback if KMS folder is missing
        
    for item in kms_prompts_dir.iterdir():
        if item.is_dir():
            name = item.name
            if "-" in name:
                parts = name.split("-", 1)
                if parts[0].isdigit():
                    max_id = max(max_id, int(parts[0]))
    return max_id + 1 if max_id != -1 else 13

def create_persona(name: str, level: str, microservice: str, specialty: str) -> None:
    workspace_root = Path(__file__).resolve().parent.parent.parent
    
    # 1. Resolve target level directory name
    obsidian_root = workspace_root / "obsidian-brain"
    level_dir = None
    for folder in obsidian_root.iterdir():
        if folder.is_dir() and folder.name.startswith(level + "-"):
            level_dir = folder
            break
            
    if not level_dir:
        # Fallback to direct name match if it is already prefixed
        if (obsidian_root / level).is_dir():
            level_dir = obsidian_root / level
        else:
            print(f"❌ Error: Could not find submodule directory for level '{level}' in obsidian-brain.")
            sys.exit(1)
            
    # 2. Determine Prefix ID from 07-Core-KMS
    kms_prompts_dir = obsidian_root / "07-Core-KMS" / "Role-Prompts"
    next_id = get_next_role_id(kms_prompts_dir)
    id_str = f"{next_id:02d}"
    
    # Format name (CamelCase for filename, PascalCase/Alphanumeric for directory)
    clean_name = "".join(x.title() for x in name.replace("-", "_").split("_"))
    folder_name = f"{id_str}-{clean_name}"
    
    role_prompts_dir = level_dir / "Role-Prompts" / folder_name
    role_prompts_dir.mkdir(parents=True, exist_ok=True)
    
    target_file = role_prompts_dir / f"Prompt-{clean_name}.md"
    if target_file.exists():
        print(f"ℹ️ Persona Prompt-{clean_name}.md already exists at {target_file.relative_to(workspace_root)}.")
        return
        
    # 3. Scaffold the template content
    template_content = f"""---
microservice: {microservice}
type: role-prompt
status: active
tags:
- '#service/{microservice}'
- '#type/role-prompt'
- '#state/active'
- '#zone/3-fleet'
---
# 🏗️ Role {id_str}: {clean_name} ({specialty})

> "A short memorable quote representing this persona's core attitude."

## 🎭 Session Initialization Ritual (MANDATORY)
You MUST begin your FIRST response in any session with the following telemetry header:
`[SCAN] Role: {clean_name} | Source: [List primary files read] | State: [Current Objective]`

## 🗂️ Context Injection (MANDATORY)
Before beginning, you MUST read:
- `Project-Variables.md` — Ecosystem constants and repo paths.
- `AI-Project-DNA.md` — Quality gates and guidelines.
- `AI-Session-State.md` — Recent session progress and tasks.

## 🎯 Primary Objective
Describe the role's primary goal and purpose within the squad.

## 🛠️ Responsibilities
1. **Core Responsibility**: [Detail primary function]
2. **Ecosystem Safety**: Enforce FFI Loading Laws and UTF-8 encoding.
3. **Session Logging**: Systematically write progress to `AI-Session-State.md` before stopping.
"""
    
    print(f"📝 Scaffolding persona folder: {role_prompts_dir.relative_to(workspace_root)}")
    with open(target_file, "w", encoding="utf-8") as f:
        f.write(template_content)
    print(f"✅ Created role prompt: {target_file.relative_to(workspace_root)}")
    
    # 4. Trigger the compiler
    convert_script = workspace_root / "obsidian-brain" / "08-Base-Scripts" / "convert_agents.py"
    if convert_script.exists():
        print("🤖 Regenerating agent persona definitions for the client squad...")
        try:
            subprocess.run([sys.executable, str(convert_script)], check=True)
            print("✨ Persona generation and registration complete!")
        except subprocess.CalledProcessError as e:
            print(f"⚠️ Failed to compile agents: {e}")
    else:
        print("⚠️ Warning: Could not find convert_agents.py script to regenerate agents.")

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Bastien-Antigravity Persona Manager")
    parser.add_argument("--name", required=True, help="Name of the persona (e.g. security-specialist)")
    parser.add_argument("--level", required=True, help="Level prefix or folder name (e.g. 03 or 03-Tech-Stack)")
    parser.add_argument("--microservice", default="core-kms-brain", help="Microservice domain")
    parser.add_argument("--specialty", required=True, help="Role specialty description")
    
    args = parser.parse_args()
    create_persona(args.name, args.level, args.microservice, args.specialty)
