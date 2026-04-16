# TODO — Improvement Suggestions

## Code Quality

- [x] **Replace panics with error returns in `internal/plan/build.go`**
  `NewBuilderFromStarConfig()` panics on failure instead of returning an error. This causes a fatal crash in the `generate` command. Refactor to return `(*Builder, error)`.

- [x] **Propagate errors in `cmd/run.go` interactive mode**
  `listCombinationDefinitions()` and `listPlanDefinitions()` silently swallow errors from `manager.List()` and return empty slices. At minimum, log errors via zerolog so failures are visible.

- [x] **Validate definition names against path traversal**
  In `cmd/definition_combination.go` and `cmd/definition_plan.go`, the `--name` flag value is used directly in file paths without validation. Add a regex check (e.g., `^[a-zA-Z0-9_-]+$`) to prevent path traversal attacks like `../../../etc/passwd`.

- [x] **Reduce code duplication in folder creation (`cmd/definition.go`)**
  The `init()` method repeats the same `createFolder` + error check pattern five times. Refactor to iterate over a slice of folder paths.

- [x] **Extract shared List() logic from service managers**
  `CombinationStarDefinitionManager.List()` and `PlanDefinitionManager.List()` have nearly identical implementations. Extract the common directory-reading logic into a shared helper or base type.

- [x] **Reduce duplication between definition_combination.go and definition_plan.go**
  Both files follow the same `New()` / `Edit()` / `Delete()` / `List()` structure. Consider a generic definition handler that works for both types.

- [x] **Make worker count configurable (`cmd/generate.go`)**
  `defaultWorkers = 2` is hardcoded with no way to override it via a CLI flag or config variable. Add an `--workers` flag or `RF_WORKERS` env variable.

- [x] **Move UI constants to a configuration struct (`cmd/run.go`)**
  Color codes (`"170"`, `"241"`, `"42"`, `"196"`), symbols, and layout constants are scattered as package-level constants. Group them into a configurable theme struct.

- [x] **Improve context propagation in `cmd/generate.go`**
  The debug server's shutdown context is created from `context.Background()` instead of inheriting from the parent context. Propagate the parent context for cleaner shutdown behavior.

## Testing

- [ ] **Add tests for `cmd/` package**
  The `cmd/` package has minimal test coverage. Critical untested areas include:
  - `run.go` — the entire interactive TUI flow (516 lines, zero tests)
  - `generate.go` — worker pool and generation orchestration
  - `definition_combination.go` / `definition_plan.go` — CRUD operations
  - `config.go` — config binding with Viper

- [ ] **Add tests for `internal/platform/fs/` package**
  No test file exists for filesystem helper functions. Cover permission errors, missing directories, and edge cases.

- [ ] **Add tests for Starlark module implementations**
  The custom Starlark modules under `internal/platform/starlark/` (random, uuid, template) have no dedicated tests. Test:
  - `random.uint()` edge cases: `min == max`, `nr > range` with `allow_duplicates=false`
  - `uuid.v7()` format validation and uniqueness
  - `template.render_text()` with complex templates and invalid input

- [ ] **Add integration tests for the full generation pipeline**
  End-to-end test: load a definition, create a plan, generate combinations, verify output files.

- [ ] **Add context cancellation tests**
  Verify that long-running operations (generation, export) properly respect `ctx.Done()` and clean up resources.

- [ ] **Add benchmark tests for plan generation**
  Benchmark the generator with varying worker counts, plan sizes, and combination complexities to identify performance bottlenecks.

## CI/CD

- [ ] **Add a GitHub Actions CI pipeline**
  No `.github/workflows/` directory exists. Create a workflow that runs on push and pull requests:
  - `go generate`
  - `go vet ./...`
  - `golangci-lint run`
  - `go test -race -cover ./...`
  - `go build`
  - Coverage report upload

- [ ] **Add a test coverage threshold**
  Coverage is output to `data/test/coverage.out` but not enforced. Add a minimum coverage gate (e.g., 70%) to the CI pipeline and a coverage badge in the README.

- [ ] **Add automated release workflow**
  Create a release workflow triggered by tags that builds cross-platform binaries and publishes a GitHub release.

## DevOps & Infrastructure

- [ ] **Optimize the production Docker image**
  The `prod` stage inherits from `base`, which includes dev tools (micro editor, Delve debugger, golangci-lint). Create a minimal final stage using `alpine` or `scratch` that contains only the compiled binary.

- [ ] **Add a `.dockerignore` file**
  Exclude `.git/`, `vendor/`, `data/`, `bin/`, and `*.md` from the Docker build context to reduce image build time and size.

- [ ] **Add health checks to `compose.yaml`**
  Neither the dev nor prod service has a health check defined. Add a simple health check to improve container orchestration reliability.

- [ ] **Add resource limits to `compose.yaml`**
  No memory or CPU limits are set. Add resource constraints to prevent runaway containers in development and production.

- [ ] **Add resource requests/limits to Kubernetes manifests**
  `k8s/job.yaml` has no resource specifications. Define CPU and memory requests/limits for predictable scheduling.

- [ ] **Replace sleep command in `k8s/job.yaml`**
  The Job currently runs `sleep 36000` instead of the actual application command. Update to use the real entrypoint.

- [ ] **Add Makefile improvements**
  - Add a `clean` target to remove generated files, binaries, and containers
  - Add `coverage` target to view HTML coverage report
  - Fix missing space in `build: generate## APP Build.` formatting
  - Validate `testPath` and `testName` variables before use

## New Features

- [ ] **Add a `--dry-run` flag to `generate combination`**
  Allow users to preview what would be generated (users, groups, file count) without writing any files. Useful for verifying plan definitions before committing to a full generation.

- [ ] **Add an `export` command for different output formats**
  Support exporting generated plans to additional formats such as PDF, HTML, or CSV for sharing with non-technical users.

- [ ] **Add a `validate` command for definitions**
  Allow users to validate Starlark and JSON definitions without generating combinations. Catch syntax errors, missing fields, and logic issues early.

- [ ] **Add plan history and diffing**
  Track previous plan generations (already stored as `.gob` files in `data/storage/`) and provide a command to compare plans across generations, showing what changed.

- [ ] **Add more Starlark built-in modules**
  Expand the scripting capabilities:
  - `math` module — common math functions (floor, ceil, abs, min, max)
  - `string` module — string manipulation utilities
  - `datetime` module — date arithmetic and formatting beyond `time.now()`
  - `regex` module — pattern matching for template-based generation

- [ ] **Support YAML plan definitions alongside JSON**
  YAML is more human-friendly for configuration. Support `.yaml` / `.yml` plan definition files in addition to `.json`.

- [ ] **Add plan templates / presets**
  Bundle common plan structures (e.g., "4-week workout split", "weekly meal plan") as built-in templates that users can start from.

## User Experience

- [ ] **Add colored output and progress bars during generation**
  The `generate combination` command runs silently. Add a progress indicator showing users generated / total combinations, elapsed time, and ETA.

- [ ] **Add tab completion for definition names**
  Register shell completions for `--combination` and `--plan` flags that auto-complete existing definition names.

- [ ] **Add a `list` command as a top-level alias**
  Instead of `random-fit definition combination` to list combination definitions, allow `random-fit list combinations` and `random-fit list plans` as more intuitive aliases.

- [ ] **Improve error messages with actionable suggestions**
  When a definition is not found or a generation fails, include suggestions like "Did you mean X?" or "Run `random-fit definition combination new --name Y` to create one."

- [ ] **Add an `init` command for first-time setup**
  Guide new users through creating their first combination definition and plan definition interactively, instead of requiring them to know the command structure.

## General

- [ ] **Add a CONTRIBUTING.md file**
  Document the development workflow, coding conventions, PR process, and how to add new Starlark modules or CLI commands.

- [ ] **Add a LICENSE file**
  The README references a LICENSE file, but none exists in the repository.

- [ ] **Add a CHANGELOG.md**
  Track notable changes across versions. Use [Keep a Changelog](https://keepachangelog.com/) format.

- [ ] **Pin dependency versions more tightly**
  Some indirect dependencies use very old versions (e.g., `gopsutil v2.19.11`). Audit and update to reduce potential supply chain risks.

- [ ] **Add Go documentation comments to all exported types and functions**
  Several exported types and functions lack doc comments. Adding them improves IDE support and enables `go doc` generation.

- [ ] **Sanitize the editor command (`cmd/definition.go`)**
  The `EDITOR` environment variable is used directly in `exec.Command()`. While mitigated by the `gosec` nolint directive, consider whitelisting known editors or at least validating the path exists and is an executable file.
