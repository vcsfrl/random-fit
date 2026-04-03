# AI Agent Context — random-fit

This file provides context for AI coding assistants (GitHub Copilot, Cursor, Claude, etc.) working on the **random-fit** project.

## Project Summary

**random-fit** is a Go CLI application that generates randomized training plans. It uses Starlark scripts to define combination logic and JSON files to define plan structures. Output is exported as JSON and Markdown files organized per user.

## Tech Stack

- **Language:** Go 1.24
- **CLI Framework:** Cobra + Viper
- **Scripting Engine:** Starlark (go.starlark.net)
- **Logging:** zerolog
- **Validation:** go-playground/validator/v10
- **Testing:** testify (suite-based)
- **Containerization:** Docker, Docker Compose, Kubernetes

## Architecture

```
main.go → cmd/ (Cobra CLI) → internal/service/ → internal/plan/ + internal/combination/
                                                   internal/platform/ (shared utilities)
```

### Key Layers

| Layer | Package | Responsibility |
|---|---|---|
| Entry point | `main.go` | Calls `cmd.Execute()` |
| CLI | `cmd/` | Cobra commands, config binding, logging |
| Service | `internal/service/` | Managers for combination and plan definitions |
| Domain | `internal/combination/` | Starlark-based combination building |
| Domain | `internal/plan/` | Plan building, generation, export |
| Platform | `internal/platform/` | Shared utilities (fs, random, starlark modules) |

### Data Flow

1. User creates a **combination definition** (`.star` file in `data/definition/`)
2. User creates a **plan definition** (`.json` file in `data/plan/`)
3. `generate combination` reads both, executes Starlark, generates combinations
4. Combinations are exported to `data/combination/{user}/{container}/{group}/`
5. Serialized plans are stored in `data/storage/` as `.gob` files

## Conventions

### Code Style

- Follow Go conventions and idioms
- Nearly all golangci-lint linters are enabled (see `.golangci.yaml`); only `exhaustruct` is disabled
- Use `zerolog` for all logging — avoid `fmt.Println` or `log.*`
- Use `go-playground/validator` for struct validation with custom validators (see `internal/combination/validator.go`)
- Errors should be wrapped with context when propagating up
- Use `context.Context` for cancellation support in long-running operations

### Naming

- Use `Star` prefix for Starlark-related types (e.g., `StarDefinition`, `StarBuilder`)
- Manager types in `internal/service/` follow the pattern `{Entity}DefinitionManager`
- Test suites use testify's `suite.Suite` and are named `{Feature}Suite`

### Testing Patterns

- **Suite-based testing** using `testify/suite` — group related tests in a struct
- **Test fixtures** in `testdata/` directories within each package
- **Mock builders** implement the `combination.Builder` interface for plan tests
- **Deterministic testing** — inject `time.Now` and `uuid.New` functions to control output in tests
- **Context cancellation tests** — verify that long-running operations respect `ctx.Done()`

### Starlark Definitions

A combination definition script must define a `definition` dict with:
- `ID` (string) — unique identifier
- `Details` (string) — human-readable description
- `BuildFunction` (function) — returns a dict of output formats

Each output format entry must have: `Extension`, `MimeType`, `Type`, `Data`.

Available built-in modules in Starlark scripts:
- `random.uint(min, max, nr, allow_duplicates, sort)`
- `uuid.v7()`
- `template.render_text(tpl, json_args)`
- `json.encode(value)` / `json.decode(str)`
- `time.now()` with `.format()` method

### Plan Definitions

JSON files with fields: `id`, `details`, `users`, `containerName`, `recurrentGroupNamePrefix`, `recurrentGroups`, `nrOfGroupCombinations`.

`containerName` supports `_date` as a special token replaced at generation time.

## Build & Development Commands

All development runs inside Docker containers via Make targets:

```bash
make install     # Build Docker image, set up .env
make shell       # Open dev shell inside container
make build       # Generate code + compile binary
make test        # Run all tests (race detection, coverage)
make lint        # Run golangci-lint (100+ linters)
make generate    # Run code generation only
make run         # Run containerized application
```

### Running a specific test

```bash
make test-name testPath=internal/plan testName=TestBuildSuite
```

### Debugging a test

```bash
make test-debug testPath=internal/plan testName=TestBuildSuite
```

## File Locations

| What | Where |
|---|---|
| CLI commands | `cmd/cmd.go` (registration), `cmd/definition_*.go`, `cmd/generate.go` |
| Config binding | `cmd/config.go` |
| Combination domain | `internal/combination/` |
| Plan domain | `internal/plan/` |
| Service managers | `internal/service/` |
| Starlark modules | `internal/platform/starlark/{random,uuid,template}/` |
| Random generator | `internal/platform/random/generator.go` |
| Starlark template | `internal/combination/template/script.star` |
| K8s manifests | `k8s/` |
| Docker config | `Dockerfile`, `compose.yaml` |
| Linter config | `.golangci.yaml` |
| Translations | `cmd/translations/` |

## Important Implementation Details

- **Combination builder** (`internal/combination/builder.go`) executes Starlark scripts in a sandboxed environment with custom modules injected as predefined names
- **Plan generator** (`internal/plan/build.go`) uses Go channels for streaming generation with worker pool support
- **Export** (`internal/plan/export.go`) writes files and serializes plans using Go's `encoding/gob`
- **Random numbers** (`internal/platform/random/generator.go`) use `crypto/rand` — never `math/rand`
- **Code generation** (`internal/service/code.go`) embeds the Starlark template file into Go source via `//go:generate`
- **The generated file** `internal/service/combination_definition_template.go` is in `.gitignore` — it is generated at build time
- **UUIDs** use UUIDv7 (time-ordered) via `google/uuid`

## Common Pitfalls

- Always run `make generate` before `make build` — the build depends on generated code (the Makefile's `build` target does this automatically)
- The `data/` directory is mounted as a Docker volume — changes inside the container are reflected on the host
- The `.env` file is gitignored; always start from `.env.dist`
- Starlark scripts use Python-like syntax but are not Python — some Python features are unavailable
- The `k8s/volumne.yaml` filename contains a deliberate typo (preserved for consistency)
