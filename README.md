# random-fit

A Go-based CLI tool for generating randomized training plans. It uses [Starlark](https://github.com/google/starlark-go) scripts (a Python-like language) to define how training combinations are built, allowing flexible and dynamic plan generation. Plans are exported in multiple formats (JSON and Markdown) for storage and distribution.

## Table of Contents

- [Features](#features)
- [Requirements](#requirements)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [CLI Commands](#cli-commands)
- [Usage Guide](#usage-guide)
- [Examples](#examples)
- [Project Structure](#project-structure)
- [Development](#development)
- [Kubernetes Deployment](#kubernetes-deployment)

## Features

- **Starlark-scripted definitions** — define combination logic using a safe, sandboxed scripting language
- **Multi-format export** — generates both JSON and Markdown output files
- **Multi-user plans** — generate personalized plans for multiple users in one run
- **Recurrent groups** — organize combinations into recurring groups (e.g., weekly schedules)
- **Cryptographically secure randomness** — uses `crypto/rand` for random number generation
- **Internationalization** — i18n support via `golang.org/x/text`
- **Docker & Kubernetes ready** — containerized development workflow and K8s deployment manifests

## Requirements

- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/)
- GNU Make

## Getting Started

### 1. Install

```bash
make install
```

This copies `.env.dist` to `.env` (if not present), appends your user ID and username, and builds the Docker image.

### 2. Open a development shell

```bash
make shell
```

This starts the dev container and opens a Bash session inside it. All `go` commands can be run directly from this shell.

### 3. Build the binary

```bash
make build
```

Generates code and compiles the binary to `bin/random-fit`.

### 4. Run the application

```bash
make run
```

Runs the containerized application using the profile set in `APP_ENV` (default: `prod`).

## Configuration

Configuration is managed via environment variables (prefix `RF_`) loaded from a `.env` file. Copy `.env.dist` as a starting point:

```bash
cp .env.dist .env
```

| Variable | Default | Description |
|---|---|---|
| `APP_ENV` | `prod` | Docker Compose profile (`prod` or `dev`) |
| `RF_DATA_FOLDER` | `/srv/random-fit/data` | Root data directory inside the container |
| `RF_BASE_FOLDER` | `/srv/random-fit` | Base application folder inside the container |
| `RF_DEBUG_CHART_PORT` | `40023` | Port for the debug metrics HTTP server |
| `RF_DEBUGGER_PORT` | `40021` | Port for Delve remote debugger |
| `RF_TRACE_PORT` | `40022` | Reserved trace port |
| `RF_LOCALE` | `en_US.UTF-8` | Locale for i18n |
| `RF_K8S_SHARED_FOLDER` | `/home/ubuntu/src/_shared_data/random-fit` | Host path for Kubernetes persistent volume |
| `EDITOR` | `micro` | Text editor used for editing definitions |

### Data Directories

All data lives under `RF_DATA_FOLDER`:

| Path | Purpose |
|---|---|
| `data/definition/` | Combination definition scripts (`.star` files) |
| `data/plan/` | Plan definition files (`.json` files) |
| `data/combination/` | Generated output (JSON and Markdown files) |
| `data/storage/` | Serialized plan objects (`.gob` binary files) |

## CLI Commands

```
random-fit
├── definition
│   ├── combination          List combination definitions
│   │   ├── new --name       Create a new combination definition
│   │   ├── edit --name      Edit an existing combination definition
│   │   └── delete --name    Delete a combination definition
│   └── plan                 List plan definitions
│       ├── new --name       Create a new plan definition
│       ├── edit --name      Edit an existing plan definition
│       └── delete --name    Delete a plan definition
├── generate
│   └── combination          Generate combinations from definitions
│       --combination        Combination definition name
│       --plan               Plan definition name
└── code
    └── generate             Generate internal template code
```

## Usage Guide

### Step 1 — Create a Combination Definition

A combination definition is a Starlark script that defines what data a single combination produces.

```bash
random-fit definition combination new --name my-workout
```

This opens your configured `$EDITOR` with a template. The script must define a `definition` dict:

```python
definition = {
    "ID": "my-workout",
    "Details": "My Workout Combination",
    "BuildFunction": build,
}
```

The `build` function must return a dict with output format keys (e.g., `json`, `markdown`). Each entry specifies:

| Key | Description |
|---|---|
| `Extension` | File extension (e.g., `json`, `md`) |
| `MimeType` | MIME type (e.g., `application/json`, `text/markdown`) |
| `Type` | Output type identifier |
| `Data` | The generated content as a string |

#### Available Starlark Modules

| Module | Function | Description |
|---|---|---|
| `random` | `random.uint(min, max, nr, allow_duplicates, sort)` | Generate random unsigned integers |
| `uuid` | `uuid.v7()` | Generate a UUIDv7 string |
| `template` | `template.render_text(tpl, json_args)` | Render a Go text template with JSON data |
| `json` | `json.encode(value)` / `json.decode(str)` | JSON serialization |
| `time` | `time.now()` | Current time (supports `.format()`) |

#### Example: Lotto Number Generator

```python
definition_id = "lotto"
definition_name = "Lotto Number Picks"

def build_combination():
    return {
        "Metadata": {
            "ID": uuid.v7(),
            "ParentID": definition_id,
            "Details": definition_name,
            "Date": time.now().format("2006-01-02T15:04:05Z07:00"),
        },
        "Data": random.uint(1, 49, 6, False, True),
    }

mdTemplate = """# {{ .Metadata.Details }}
##### Date: {{ .Metadata.Date }}
[ {{ range .Data }}{{.}} {{ end }}]
"""

def build():
    combination = build_combination()
    json_combination = json.encode(combination)
    return {
        "json": {
            "Extension": "json",
            "MimeType": "application/json",
            "Type": "json",
            "Data": json_combination,
        },
        "markdown": {
            "Extension": "md",
            "MimeType": "text/markdown",
            "Type": "markdown",
            "Data": template.render_text(mdTemplate, json_combination),
        },
    }

definition = {
    "ID": definition_id,
    "Details": definition_name,
    "BuildFunction": build,
}
```

### Step 2 — Create a Plan Definition

A plan definition is a JSON file that describes users, groups, and how many combinations to generate.

```bash
random-fit definition plan new --name my-plan
```

This opens your editor with a JSON template. Example:

```json
{
  "id": "gym_plan",
  "details": "Weekly Gym Training",
  "users": ["alice", "bob"],
  "containerName": ["_date", "gym"],
  "recurrentGroupNamePrefix": "Week",
  "recurrentGroups": 4,
  "nrOfGroupCombinations": 3
}
```

| Field | Description |
|---|---|
| `id` | Unique identifier for the plan |
| `details` | Human-readable description |
| `users` | List of user names to generate plans for |
| `containerName` | Folder path segments (`_date` is replaced with a timestamp) |
| `recurrentGroupNamePrefix` | Prefix for group folders (e.g., `Week-1`, `Week-2`) |
| `recurrentGroups` | Number of recurring groups to generate |
| `nrOfGroupCombinations` | Number of combinations per group |

### Step 3 — Generate Combinations

```bash
random-fit generate combination --combination my-workout --plan my-plan
```

Or using positional arguments:

```bash
random-fit generate combination my-workout my-plan
```

This produces output in `data/combination/` organized by user, date, container, and group:

```
data/combination/
├── alice/
│   └── 2025-01-15-10-30/gym/
│       ├── Week-1/
│       │   ├── My-Workout_1.json
│       │   ├── My-Workout_1.md
│       │   ├── My-Workout_2.json
│       │   └── My-Workout_2.md
│       ├── Week-2/
│       │   └── ...
│       └── ...
└── bob/
    └── ...
```

### Managing Definitions

```bash
# List all combination definitions
random-fit definition combination

# Edit an existing combination definition
random-fit definition combination edit --name my-workout

# Delete a combination definition
random-fit definition combination delete --name my-workout

# List all plan definitions
random-fit definition plan

# Edit / delete plan definitions
random-fit definition plan edit --name my-plan
random-fit definition plan delete --name my-plan
```

## Examples

The [`examples/`](examples/) directory contains three ready-to-use use cases
that you can copy into your `data/` folders and run immediately.

| Use Case | Combination definition | Plan definition | What it generates |
|---|---|---|---|
| **Gym Workout** | `examples/definition/gym-workout.star` | `examples/plan/gym-workout-plan.json` | 4-week programme, 3 sessions/week, 2 users — each session: 3 compound + 2 isolation exercises with random sets/reps |
| **Lotto 6/49** | `examples/definition/lotto.star` | `examples/plan/lotto-plan.json` | 4-week picks, 5 tickets/week, 3 users — each ticket: 6 sorted main numbers (1–49) + 1 bonus number (1–10) |
| **Meal Plan** | `examples/definition/meal-plan.star` | `examples/plan/meal-plan.json` | 4-week daily plan, 7 days/week, 2 users — each day: breakfast, lunch, dinner, snack + estimated calorie total |

### Quick start with an example

```bash
# Copy a definition and its plan into the data directories
cp examples/definition/gym-workout.star data/definition/
cp examples/plan/gym-workout-plan.json  data/plan/

# Generate — output lands in data/combination/<user>/<timestamp>/gym/
random-fit generate combination --combination gym-workout --plan gym-workout-plan
```

See [`examples/README.md`](examples/README.md) for detailed explanations, sample output, and the full output directory structure for each use case.

## Project Structure

```
random-fit/
├── main.go                          # Entry point
├── cmd/                             # CLI commands (Cobra)
│   ├── cmd.go                       # Root command and subcommand registration
│   ├── config.go                    # Configuration binding (Viper)
│   ├── definition.go                # Definition command group
│   ├── definition_combination.go    # Combination CRUD handlers
│   ├── definition_plan.go           # Plan CRUD handlers
│   ├── generate.go                  # Generate command handler
│   ├── log.go                       # Logger setup (zerolog)
│   └── translations/                # i18n catalogs (gotext)
├── internal/
│   ├── combination/                 # Combination domain logic
│   │   ├── combination.go           # Core data structures
│   │   ├── definition.go            # Starlark script loader
│   │   ├── builder.go               # Combination builder
│   │   ├── validator.go             # Custom validators
│   │   └── template/script.star     # Default Starlark template
│   ├── plan/                        # Plan domain logic
│   │   ├── plan.go                  # Plan/UserPlan structs, JSON loader
│   │   ├── build.go                 # Plan builder and generator
│   │   └── export.go                # File exporter and GOB serializer
│   ├── service/                     # Service layer
│   │   ├── combination.go           # CombinationStarDefinitionManager
│   │   ├── plan.go                  # PlanDefinitionManager
│   │   └── config.go                # Folder configuration
│   └── platform/                    # Shared utilities
│       ├── fs/                      # Filesystem helpers
│       ├── random/                  # Crypto-based random generator
│       └── starlark/                # Custom Starlark modules
│           ├── random/              # random.uint()
│           ├── uuid/                # uuid.v7()
│           └── template/            # template.render_text()
├── data/                            # Runtime data (mounted volume)
├── k8s/                             # Kubernetes manifests
├── Dockerfile                       # Multi-stage Docker build
├── Makefile                         # Build targets
├── compose.yaml                     # Docker Compose (dev + prod)
├── .env.dist                        # Environment template
└── .golangci.yaml                   # Linter configuration
```

## Development

### Prerequisites

All development is done inside a Docker container. Start the dev shell:

```bash
make shell
```

### Make Targets

| Target | Description |
|---|---|
| `make help` | Show all available targets |
| `make install` | Build Docker image and create `.env` |
| `make shell` | Open a Bash shell in the dev container |
| `make generate` | Run code generation (templates + translations) |
| `make build` | Generate code and compile the binary |
| `make test` | Run all tests with race detection and coverage |
| `make test-name testPath=<pkg> testName=<name>` | Run a single test by name |
| `make test-debug testPath=<pkg> testName=<name>` | Debug a test with Delve |
| `make lint` | Run golangci-lint (100+ linters enabled) |
| `make run` | Run the application in a container |

### Testing

Tests use [testify](https://github.com/stretchr/testify) with suite-based organization. Run the full suite:

```bash
make test
```

Run a specific test:

```bash
make test-name testPath=internal/plan testName=TestBuildSuite
```

Coverage output is written to `data/test/coverage.out`.

### Linting

The project uses [golangci-lint v2](https://golangci-lint.run/) with nearly all available linters enabled (see `.golangci.yaml`):

```bash
make lint
```

### Key Dependencies

| Library | Purpose |
|---|---|
| [cobra](https://github.com/spf13/cobra) | CLI framework |
| [viper](https://github.com/spf13/viper) | Configuration management |
| [starlark-go](https://github.com/google/starlark-go) | Starlark script execution |
| [zerolog](https://github.com/rs/zerolog) | Structured logging |
| [validator](https://github.com/go-playground/validator) | Data validation |
| [uuid](https://github.com/google/uuid) | UUIDv7 generation |
| [testify](https://github.com/stretchr/testify) | Testing utilities |
| [statsviz](https://github.com/arl/statsviz) | Runtime debug metrics |

## Kubernetes Deployment

The `k8s/` directory contains manifests for deploying as a Kubernetes Job (tested with MicroK8s).

### Deploy

```bash
cd k8s
./install.sh
```

This creates a ConfigMap, a PersistentVolumeClaim (1Gi), and a Job that runs the application.

### Remove

```bash
cd k8s
./uninstall.sh
```

### Manifests

| File | Resource |
|---|---|
| `job-config.yaml` | ConfigMap with environment settings |
| `volumne.yaml` | StorageClass + PersistentVolumeClaim (hostpath) |
| `job.yaml` | Kubernetes Job running the `random-fit` container |

## License

See [LICENSE](LICENSE) for details.
