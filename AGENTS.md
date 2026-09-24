---
microservice: 08-Base-Scripts
type: note
status: active
tags:
- '#service/08-Base-Scripts'
- '#type/note'
- '#state/active'
- '#zone/3-fleet'
---# AGENTS.md: nexus-strategic-brain (01-Strategic-Nexus)

## Service Mission & Architecture Role
`01-Strategic-Nexus` is the strategic orchestration, evolutionary governance, and milestone management microservice for the Bastien-Antigravity fleet. It manages strategic intent, connects with the TimescaleDB (`obsidiandb`), runs AI prompt evolutionary pipelines via Gemini, hosts an OpenMFE micro-frontend, and provides remote interactive controls over Telegram via `tele-remote`.

- **Exposed Capability**: `strategic_nexus`
  - **OpenMFE REST**: Resolved dynamically via `appConfig.GetCapability("strategic_nexus", &struct)` or `capabilities.strategic_nexus.mfe_port` (schema default: `8092`)
- **Downstream Integrations**:
  - `tele-remote`: Resolved via `appConfig.GetGRPCAddr("tele_remote")`
  - `timescale_db`: Resolved via `appConfig.GetCapability("timescale_db", &struct)`
  - `web-interface`: Resolved via `appConfig.GetCapability("web_interface", &struct)` for OpenMFE dynamic registration
- **Shared Libraries**: `microservice-toolbox`, `universal-logger`, `distributed-config`
- **Configuration Link**: `standalone.yaml -> ../../docker-deployment/modes/local/config/native.yaml`

> [!CAUTION]
> **CRITICAL ANTI-PATTERN WARNING FOR AI AGENTS**:
> NEVER hardcode IP addresses (`127.0.0.1`, `0.0.0.0`, `localhost`) or port numbers in source code. All network endpoints, hostnames, ports, and database credentials MUST be resolved dynamically at runtime through `microservice-toolbox` configuration accessors. Port numbers mentioned in documentation represent ONLY canonical defaults configured in `native.yaml`.

## Key Build & Execution Commands
```bash
# Display service version
make version

# Build binary
make build
# or directly:
go build -o bin/strategic-nexus ./cmd/strategic-nexus

# Run unit tests
make test
# or directly:
go test -v ./...

# Run service daemon
make run
# or directly:
./bin/strategic-nexus

# Run direct CLI commands
./bin/strategic-nexus run-evolution --prompt-name <name>
```

## Internal Architecture & Subsystems
1. **Core Controller (`src/core/`)**: Coordinates evolutionary pipelines, milestone logging, and state transitions.
2. **Evolution Pipeline (`src/evolution/`)**: Integrates with Gemini API to evaluate and evolve agent prompts and strategic directives.
3. **Storage & Migrations (`src/store/`)**: Connects to `obsidiandb` using `pgxpool` with automated schema migrations.
4. **REST & OpenMFE Host (`src/rest/`)**: Hosts the Strategic Nexus OpenMFE micro-frontend and dynamic registration with `web-interface`.
5. **Telegram Integration (`src/telegram/`)**: Subscribes to `tele-remote` for command and telemetry dispatch.

## AI Development & Integration Guidelines
1. **Triple-Block Header**: All Go source files MUST begin with the Triple-Block header (`ESSENTIAL PROCESS`, `DATA FLOW`, `KEY PARAMETERS`).
2. **Dynamic Capability Resolution**: Resolve addresses via `appConfig.GetCapability(...)`, never hardcode host or port values.
3. **Encrypted Credentials**: Database passwords and API tokens must be decrypted using `appConfig.DecryptSecret()`.
4. **Lifecycle Registration**: All servers must be registered with `toolbox_lifecycle.Manager` for graceful shutdown on `SIGINT`/`SIGTERM`.
