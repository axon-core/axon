# Project Conventions for AI Assistants

## Rules for AI Assistants
- **Use Makefile targets** instead of discovering build/test commands yourself.
- **Keep changes minimal.** Do not refactor, reorganize, or 'improve' code beyond what was explicitly requested.
- **For CI/release workflows**, always use existing Makefile targets rather than reimplementing build logic in YAML.
- **Better tests.** Always try to add or improve tests(including integration, e2e) when modifying code.
- **Logging conventions.** Start log messages with capital letters and do not end with punctuation.
- **Commit messages.** Do not include PR links in commit messages.
- **Pin container image tags.** Always use specific version tags (e.g., `node:22.14.0-alpine`), never rolling tags (e.g., `node:22-alpine`) to ensure reproducible builds.
- **GitHub Actions permissions.** Always add an explicit `permissions` block to workflow files, especially those using `pull_request_target`. Follow least-privilege: only request the permissions the workflow actually needs, or use `permissions: {}` if none are needed.
- **Use typed constants in Go tests.** Use typed constants from packages (e.g., `batchv1.PodFailurePolicyActionCount`) instead of string literals (e.g., `"Count"`) when comparing against values that have defined constants.

## Key Makefile Targets
- `make verify` — run all verification checks (lint, fmt, vet, etc.).
- `make update` — update all generated files
- tests:
  - `make test` — run all unit tests
  - `make test-integration` — run integration tests
  - e2e tests are hard to run locally. Push changes and use the PR's CI jobs to run them instead.
- `make build` — build binary

## Pull Requests
- **Always follow `.github/PULL_REQUEST_TEMPLATE.md`** when creating PRs.
- Fill in every section of the template. Do not remove or skip sections — use "N/A" or "NONE" where appropriate.
- Choose exactly one `/kind` label from: `bug`, `cleanup`, `docs`, `feature`.
- If there is no associated issue, write "N/A" under the issue section.
- If the PR does not introduce a user-facing change, write "NONE" in the `release-note` block.

## Directory Structure
- `cmd/` — CLI entrypoints
- `test/e2e/` — end-to-end tests
- `.github/workflows/` — CI workflows

