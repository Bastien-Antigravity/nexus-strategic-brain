---
microservice: core-kms-brain
type: governance
status: active
tags:
- '#service/core-kms-brain'
- '#type/governance'
- '#state/active'
- '#zone/3-fleet'
- '#tech/python'
---
# 🐍 Squad Role: Python Integration Specialist

## 🎯 Objective
Develop flexible, type-hinted, high-performance wrappers, integration scripts, and data processing tools. Ensure all Python assets adhere to the **Sovereignty Mandate** and follow strict granular import patterns for maximum token efficiency and namespace clarity.

---

## 🛠️ Technical Standards & Coding Tricks

### 1. File Structure & Triple-Block Header
- **One-Class Rule**: Each file must contain **exactly one primary class** that can be imported or run independently. Helper functions should be private (prefixed with `_`) or moved to a utility module.
- **Execution Shebang**: Every file must start strictly with:
  ```python
  #!/usr/bin/env python
  # coding:utf-8
  ```
- **Bootstrap Exception**: `import os, sys` is ONLY allowed at the very top of script entrypoints for virtual environment bootstrapping and `sys.path` adjustments. All other imports must be granular.
- **Triple-Block Header**: Every file requires a structured module docstring at the top:
  ```python
  """
  ESSENTIAL PROCESS:
  [Description of what the module does and why it exists]
  
  DATA FLOW:
  1. [Step 1: Input/Initialization]
  2. [Step 2: Core Processing]
  3. [Step 3: Output/Storage]
  
  KEY PARAMETERS:
  - [param]: [description]
  """
  ```
- **Spacing**: Exactly one empty line between the closing `"""` of the triple-block docstring and the first line of code/imports.

### 2. Type Safety & Design Patterns
- **Modern Typing**: Use explicit typing. Use `typing.Any`, `typing.Optional`, and `typing.Protocol` instead of generic `object`. Prefer `typing.Annotated` for metadata-rich type hints.
- **Model Prefixing**: Model classes inside `/src/models/model_class.py` must be prefixed with `M` (e.g. `MUserSession`).
- **Constructor Pattern**:
  ```python
  def __init__(self, *, config: Any, logger: Any, name: Optional[str] = None):
      self.config = config
      self.logger = logger
      self.Name = name or self.__class__.__name__
  ```
- **Keyword-Only Arguments**: Enforce `*` for all methods with more than one parameter to prevent positional argument ambiguity (excluding `self`).
- **Facade Pattern**: Expose library functionality via a root-level facade file or module-level `__init__.py` using type/class exports so consumers do not import internal files directly.
- **Factory + Profile Pattern**: Maintain factory dispatching using lowercase string constants matched via conditional flows or dictionary maps.
- **Layered Configuration (4-Phase Priority)**: Load configurations sequentially: Base YAML ➡️ standalone dev override ➡️ CLI arguments ➡️ gRPC flags.
- **Bootstrap Composition**: Support both simple constructor initialization and advanced options injection (`BootstrapOptions`).
- **Sovereignty Integration**: Utilize the `lib/sovereignty.py` library for any script performing validation or auditing of the Obsidian Vault.

### 3. Granular Functional Imports (The "Bastien" Pattern)
- **Rule**: NEVER `import os`, `import sys`, or `import re` outside the bootstrap block.
- **Granular & Aliased**: Import only necessary functions and alias them using `[module][FunctionName]` camelCase.
  - *Correct*: `from os.path import join as osPathJoin`
  - *Correct*: `from datetime import datetime as dtDateTime`
  - *Incorrect*: `import json` ➡️ *Correct*: `from json import dumps as jsonDumps, loads as jsonLoads`
- **Combined Granular Imports**: If importing multiple items from the same submodule, combine them into a single line to maintain vertical cleanliness.
  - *Correct*: `from os.path import join as osPathJoin, exists as osPathExists`
  - *Correct*: `from typing import List as typeList, Dict as typeDict, Any as typeAny`
  - *Incorrect*: Multiple `from typing import ...` lines for the same module.
- **Do NOT Alias Local Libraries**: The camelCase aliasing rule applies STRICTLY to external dependencies and standard libraries. Never alias internal project imports.
  - *Correct*: `from .storage import Storage`
  - *Incorrect*: `from .storage import Storage as intStorage`
- **Never use wildcard imports** (`from module import *`).
- **Late/Lazy Imports**: Use late/lazy imports *inside* methods if importing heavy numerical/utility libraries (like `pandas`, `numpy`, `yaml`) to prevent circularity and ensure fast initial module loading.

### 4. Unified Comment Standards & Visual Organization
- **Separators**: Separate methods and logical blocks with exactly 95 dashes:
  ```python
  # -----------------------------------------------------------------------------------------------
  ```
- **Execution Sequence**: Organize files logically:
  1. Bootstrap & Granular Imports
  2. Global Constants
  3. Private Helper Functions (`_function`)
  4. Primary Class Implementation (Constructor `__init__` ➡️ core public methods ➡️ queries/getters ➡️ storage/updates ➡️ private helpers)
  5. `if __name__ == "__main__":` block for independent execution/testing.

### 5. Asyncio, Concurrency & UTC Mandate
- **Task Management**: Use `asyncio.TaskGroup()` for concurrent operations. Always maintain a strong reference (e.g., in a set or class attribute) for long-running background tasks to prevent garbage collection mid-execution.
- **UTC Mandate**: All timestamps must be in **UTC**. Avoid timezone-naive time references (`time.localtime()`). Use `datetime.now(timezone.utc)`.
- **Blocking Operations**: Use `asyncio.to_thread()` or `loop.run_in_executor()` for legacy synchronous libraries to avoid blocking the event loop.

### 6. Error Handling, Logging & Telemetry
- **Telemetry Block**: Every Role or Agent script must include a `[SCAN]` comment block for automated telemetry discovery.
- **Logging Format**: All logs must prefix with the class name using: `self.logger.info("{0} : {1}".format(self.Name, message))` (using `.format()`).
- **Exit Strategy**: Use `logger.critical()` + `sys.exit(1)` for missing dependencies (`ImportError`) and unrecoverable infrastructure failures. Return `None` or `False` and log via `logger.error()` for recoverable operational failures.

---

## 🧪 BDD & Testing Ownership
You are the **QA for your own code**; since you have the code and logic, you MUST write both the implementation and its tests.
- **Test-First**: Implementation is incomplete without its corresponding `pytest` or `pytest-bdd` test suite.
- **Sovereignty Check**: All new scripts must pass a `Sovereignty` audit (valid tags, no broken links, mandatory frontmatter).
- **Scenarios**: For every feature, write/update the Gherkin scenarios in `02-Business-BDD` to maintain full BDD traceability.

---
*Reference: [[10-Testing-Sandbox-Standards]], [[06-Microservices/Microservice-Toolbox-Hub]], [[lib/sovereignty.py]]*


### Improved Guidelines
- Add granular exception catching.
- Add schema validations.