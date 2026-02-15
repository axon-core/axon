# Axon

**The Kubernetes-native framework for orchestrating autonomous AI coding agents.**

[![CI](https://github.com/axon-core/axon/actions/workflows/ci.yaml/badge.svg)](https://github.com/axon-core/axon/actions/workflows/ci.yaml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/axon-core/axon)](https://github.com/axon-core/axon)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![GitHub Release](https://img.shields.io/github/v/release/axon-core/axon)](https://github.com/axon-core/axon/releases/latest)

- Run AI coding agents (Claude Code, Codex, Gemini, OpenCode) as autonomous Kubernetes workloads
- Chain tasks into pipelines, react to GitHub events, and fan out across repos
- Every agent runs in an isolated, ephemeral Pod — no access to your host machine

<p align="center">
  <video src="https://github.com/user-attachments/assets/837cd8d5-4071-42dd-be32-114c649386ff" width="800" controls></video>
</p>

```bash
axon init                                          # configure credentials
axon run -p "Fix issue #42 and open a PR"          # launch an agent
axon logs task-a5b3c -f                            # stream the output
```

## Quick Start

### Prerequisites

- Kubernetes cluster (1.28+) — don't have one? `kind create cluster`
- kubectl configured

### 1. Install the CLI

```bash
curl -fsSL https://raw.githubusercontent.com/axon-core/axon/main/hack/install.sh | bash
```

Or from source: `go install github.com/axon-core/axon/cmd/axon@latest`

### 2. Install Axon

```bash
axon install
```

This installs the Axon controller and CRDs into the `axon-system` namespace.

### 3. Initialize and run

```bash
axon init
# Edit ~/.axon/config.yaml with your token and workspace:
#   oauthToken: <your-oauth-token>
#   workspace:
#     repo: https://github.com/your-org/your-repo.git
#     ref: main
#     token: <github-token>  # optional, for private repos and pushing changes
```

<details>
<summary>How to get your credentials</summary>

**Claude OAuth token** (recommended for Claude Code):
Run `claude auth login` locally, then copy the token from `~/.claude/credentials.json`.

**Anthropic API key** (alternative):
Create one at [console.anthropic.com](https://console.anthropic.com). Set `apiKey` instead of `oauthToken` in your config.

**GitHub token** (for pushing branches and creating PRs):
Create a [Personal Access Token](https://github.com/settings/tokens) with `repo` scope (and `workflow` if your repo uses GitHub Actions).

</details>

### 4. Run your first task

```bash
$ axon run -p "Add a hello world program in Python"
task/task-r8x2q created

$ axon logs task-r8x2q -f
```

The agent clones your repo, makes changes, and can push a branch or open a PR. Use `--name` for a custom name, or `-w` to auto-watch logs.

> **Tip:** If something goes wrong, check controller logs with
> `kubectl logs deployment/axon-controller-manager -n axon-system`.

## Why Axon?

- **Orchestration, not just execution** — Chain tasks with `dependsOn`, pass results between stages, and use `TaskSpawner` to build event-driven workers that react to GitHub issues, PRs, or schedules.
- **Host-isolated autonomy** — Each task runs in an isolated, ephemeral Pod with a freshly cloned workspace. Agents have no access to your host machine.
- **Standardized interface** — Plug in any agent (Claude, Codex, Gemini, OpenCode, or your own) using a simple [container interface](docs/agent-image-interface.md). Axon handles credentials, workspaces, and Kubernetes plumbing.
- **Scalable parallelism** — Fan out agents across repositories. Kubernetes handles scheduling and resource management — scale is limited only by cluster capacity and API quotas.
- **Observable and CI-native** — Every run is a Kubernetes resource with deterministic outputs (branch, PR URL, commit SHA, token usage). Monitor via `kubectl`, manage via CLI or declarative YAML (GitOps-ready).

## How It Works

```
  Triggers (GitHub, Cron) ──┐
                            │
  Manual (CLI, YAML) ───────┼──▶  TaskSpawner  ──▶  Tasks  ──▶  Isolated Pods
                            │          │              │             │
  API (CI/CD, Webhooks) ────┘          └─(Lifecycle)──┴─(Execution)─┴─(Success/Fail)
```

You define what needs to be done, and Axon handles the "how" — from cloning the right repo and injecting credentials to running the agent and capturing its outputs (branch names, commit SHAs, PR URLs, and token usage).

<details>
<summary>TaskSpawner — Automatic Task Creation from External Sources</summary>

TaskSpawner watches external sources (e.g., GitHub Issues) and automatically creates Tasks for each discovered item.

```
                    polls         new issues
 TaskSpawner ─────────────▶ GitHub Issues
      │        ◀─────────────
      │
      ├──creates──▶ Task: fix-bugs-1
      └──creates──▶ Task: fix-bugs-2
```

</details>

## Who Is This For?

- **Platform engineers** building AI-assisted development workflows on Kubernetes
- **DevOps teams** automating code changes, migrations, and refactoring at scale
- **Engineering teams** wanting to run AI agents safely in isolated, ephemeral environments

## Examples

### Run against a git repo

```yaml
# ~/.axon/config.yaml
oauthToken: <your-oauth-token>
workspace:
  repo: https://github.com/your-org/repo.git
  ref: main
```

```bash
axon run -p "Add unit tests"
```

Or reference an existing Workspace resource: `axon run -p "Add unit tests" --workspace my-workspace`

### Create PRs automatically

Add a `token` to your workspace config:

```yaml
workspace:
  repo: https://github.com/your-org/repo.git
  ref: main
  token: <your-github-token>
```

```bash
axon run -p "Fix the bug described in issue #42 and open a PR with the fix"
```

The `gh` CLI and `GITHUB_TOKEN` are available inside the agent container, so the agent can push branches and create PRs autonomously.

### Inject agent instructions, plugins, and MCP servers

Use `AgentConfig` to bundle project-wide instructions (like `AGENTS.md` or `CLAUDE.md`), Claude Code plugins (skills and agents), and MCP servers. Tasks reference it via `agentConfigRef`:

```yaml
apiVersion: axon.io/v1alpha1
kind: AgentConfig
metadata:
  name: my-config
spec:
  agentsMD: |
    # Project Rules
    Follow TDD. Always write tests first.
  plugins:
    - name: team-tools
      skills:
        - name: deploy
          content: |
            ---
            name: deploy
            description: Deploy the application
            ---
            Deploy instructions here...
      agents:
        - name: reviewer
          content: |
            ---
            name: reviewer
            description: Code review specialist
            ---
            You are a code reviewer...
  mcpServers:
    - name: github
      type: http
      url: https://api.githubcopilot.com/mcp/
      headers:
        Authorization: "Bearer <token>"
    - name: my-tools
      type: stdio
      command: npx
      args: ["-y", "@my-org/mcp-tools"]
```

Reference from a Task via YAML (`spec.agentConfigRef.name: my-config`) or the CLI:

```bash
axon create agentconfig my-config \
  --agents-md @AGENTS.md \
  --skill deploy=@skills/deploy.md \
  --agent reviewer=@agents/reviewer.md \
  --mcp 'github={"type":"http","url":"https://api.githubcopilot.com/mcp/"}' \
  --mcp 'my-tools=@mcp/my-tools.json'

axon run -p "Fix the bug" --agent-config my-config
```

- `agentsMD` is written to `~/.claude/CLAUDE.md` (user-level, additive with repo instructions).
- `plugins` are mounted as plugin directories and passed via `--plugin-dir`.
- `mcpServers` are written to the agent's native MCP configuration (e.g., `~/.claude.json` for Claude Code, `~/.codex/config.toml` for Codex, `~/.gemini/settings.json` for Gemini). Supports `stdio`, `http`, and `sse` transport types.

### Auto-fix GitHub issues with TaskSpawner

```yaml
apiVersion: axon.io/v1alpha1
kind: TaskSpawner
metadata:
  name: fix-bugs
spec:
  when:
    githubIssues:
      labels: [bug]
      state: open
  taskTemplate:
    type: claude-code
    workspaceRef:
      name: my-workspace
    credentials:
      type: oauth
      secretRef:
        name: claude-oauth-token
    promptTemplate: "Fix: {{.Title}}\n{{.Body}}"
  pollInterval: 5m
```

TaskSpawner polls for new issues matching your filters and creates a Task for each one. You can also run tasks on a cron schedule — see the [examples/](examples/) directory.

### Chain tasks with dependencies

Use `dependsOn` to chain tasks into pipelines. A task in `Waiting` phase stays paused until all its dependencies succeed:

```bash
axon run -p "Scaffold a new user service" --name scaffold --branch feature/user-service
axon run -p "Write tests for the user service" --depends-on scaffold --branch feature/user-service
```

Tasks sharing the same `branch` are serialized automatically — only one runs at a time.

<details>
<summary>YAML equivalent</summary>

```yaml
apiVersion: axon.io/v1alpha1
kind: Task
metadata:
  name: scaffold
spec:
  type: claude-code
  prompt: "Scaffold a new user service with CRUD endpoints"
  credentials:
    type: oauth
    secretRef:
      name: claude-oauth-token
  workspaceRef:
    name: my-workspace
  branch: feature/user-service
---
apiVersion: axon.io/v1alpha1
kind: Task
metadata:
  name: write-tests
spec:
  type: claude-code
  prompt: "Write comprehensive tests for the user service"
  credentials:
    type: oauth
    secretRef:
      name: claude-oauth-token
  workspaceRef:
    name: my-workspace
  branch: feature/user-service
  dependsOn: [scaffold]
```

</details>

### Autonomous self-development pipeline

This is a real-world TaskSpawner that picks up every open issue, investigates it, opens (or updates) a PR, self-reviews, and ensures CI passes — fully autonomously. When the agent can't make progress, it labels the issue `axon/needs-input` and stops. Remove the label to re-queue it.

```
 ┌──────────────────────────────────────────────────────────────────┐
 │                        Feedback Loop                             │
 │                                                                  │
 │  ┌─────────────┐  polls  ┌────────────────┐                     │
 │  │ TaskSpawner │───────▶ │ GitHub Issues  │                     │
 │  └──────┬──────┘         │ (open, no      │                     │
 │         │                │  needs-input)  │                     │
 │         │ creates        └────────────────┘                     │
 │         ▼                                                       │
 │  ┌─────────────┐  runs   ┌─────────────┐  opens PR   ┌───────┐ │
 │  │    Task     │───────▶ │    Agent    │────────────▶│ Human │ │
 │  └─────────────┘  in Pod │   (Claude)  │  or labels  │Review │ │
 │                          └─────────────┘  needs-input└───┬───┘ │
 │                                                          │     │
 │                                           removes label ─┘     │
 │                                           (re-queues issue)    │
 └────────────────────────────────────────────────────────────────┘
```

See [`self-development/axon-workers.yaml`](self-development/axon-workers.yaml) for the full manifest and the [`self-development/` README](self-development/README.md) for setup instructions.

The key pattern is `excludeLabels: [axon/needs-input]` — this creates a feedback loop where the agent works autonomously until it needs human input, then pauses. Removing the label re-queues the issue on the next poll.

More examples in the [`examples/`](examples/) directory — ready-to-apply YAML manifests for common use cases.

## Orchestration Patterns

- **Autonomous Self-Development** — Build a feedback loop where agents pick up issues, write code, self-review, and fix CI flakes until the task is complete.
- **Event-Driven Bug Fixing** — Automatically spawn agents to investigate and fix bugs as soon as they are labeled in GitHub.
- **Fleet-Wide Refactoring** — Orchestrate a "fan-out" where dozens of agents apply the same refactoring pattern across a fleet of microservices in parallel.
- **Hands-Free CI/CD** — Embed agents as first-class steps in your deployment pipelines to generate documentation or perform automated migrations.
- **AI Worker Pools** — Maintain a pool of specialized agents (e.g., "The Security Fixer") that developers can trigger via simple Kubernetes resources.

## Reference

Full specification for all Axon resources, configuration, and CLI commands:

**[docs/reference.md](docs/reference.md)** — Task Spec, Workspace Spec, AgentConfig Spec, TaskSpawner Spec, promptTemplate Variables, Task & TaskSpawner Status, Configuration, and CLI Reference.

<details>
<summary>Using kubectl and YAML instead of the CLI</summary>

Create a `Workspace` resource to define a git repository:

```yaml
apiVersion: axon.io/v1alpha1
kind: Workspace
metadata:
  name: my-workspace
spec:
  repo: https://github.com/your-org/your-repo.git
  ref: main
```

Then reference it from a `Task`:

```yaml
apiVersion: axon.io/v1alpha1
kind: Task
metadata:
  name: hello-world
spec:
  type: claude-code
  prompt: "Create a hello world program in Python"
  credentials:
    type: oauth
    secretRef:
      name: claude-oauth-token
  workspaceRef:
    name: my-workspace
```

```bash
kubectl apply -f workspace.yaml
kubectl apply -f task.yaml
kubectl get tasks -w
```

</details>

<details>
<summary>Using an API key instead of OAuth</summary>

Set `apiKey` instead of `oauthToken` in `~/.axon/config.yaml`:

```yaml
apiKey: <your-api-key>
```

Or pass `--secret` to `axon run` with a pre-created secret (api-key is the default credential type), or set `spec.credentials.type: api-key` in YAML.

</details>

## Roadmap

- Multi-cluster support
- Web dashboard / UI
- Agent result aggregation and reporting
- Webhook triggers (in addition to polling)
- Cost tracking and budgets per team/project

## Security Considerations

Axon runs agents in isolated, ephemeral Pods — they cannot access your host machine, SSH keys, or other processes. However, agents **do** have write access to your repositories and GitHub API via injected credentials. Here's how to manage the risk:

- **Scope your GitHub tokens.** Use [fine-grained Personal Access Tokens](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#fine-grained-personal-access-tokens) restricted to specific repositories.
- **Enable branch protection.** Require PR reviews before merging to `main`. Agents can push branches and open PRs, but protected branches prevent direct pushes.
- **Use `maxConcurrency` and `maxTotalTasks`.** Limit how many tasks a TaskSpawner can create to prevent runaway activity.
- **Use `podOverrides.activeDeadlineSeconds`.** Set a timeout to prevent tasks from running indefinitely.
- **Audit via Kubernetes.** Every agent run is a first-class Kubernetes resource — use `kubectl get tasks` and cluster audit logs.

> **Why `--dangerously-skip-permissions`?** Claude Code uses this flag to run without interactive approval prompts, which is necessary for autonomous execution. In Axon's context the agent runs inside an ephemeral container with no host access — the flag allows non-interactive operation. The actual risk surface is limited to what the injected credentials allow.

Axon uses standard Kubernetes RBAC — use namespace isolation to separate teams. Each TaskSpawner automatically creates a scoped ServiceAccount and RoleBinding.

## Cost and Limits

Running AI agents costs real money. Here's how to stay in control:

**Model costs vary significantly.** Opus is the most capable but most expensive model. Use `spec.model` (or `model` in config) to choose cheaper models like Sonnet for routine tasks and reserve Opus for complex work.

**Use `maxConcurrency` to cap spend.** Without it, a TaskSpawner can create unlimited concurrent tasks. Always set a limit:

```yaml
spec:
  maxConcurrency: 3      # max 3 tasks running at once
  maxTotalTasks: 50       # stop after 50 total tasks
```

**Use `podOverrides.activeDeadlineSeconds` to limit runtime:**

```bash
axon run -p "Fix the bug" --timeout 30m
```

**Use `suspend` for emergencies:**

```bash
axon suspend taskspawner my-spawner
# ... investigate ...
axon resume taskspawner my-spawner
```

**Rate limits.** API providers enforce concurrency and token limits. Use `maxConcurrency` to stay within your provider's limits.

## Community

We're building Axon in the open. Get involved:

- [GitHub Issues](https://github.com/axon-core/axon/issues) — bug reports and feature requests
- [GitHub Discussions](https://github.com/axon-core/axon/discussions) — questions, ideas, show & tell

If you find Axon useful, consider giving us a star — it helps others discover the project.

## Development

```bash
make update             # generate code, CRDs, fmt, tidy
make verify             # generate + vet + tidy-diff check
make test               # unit tests
make test-integration   # integration tests (envtest)
make test-e2e           # e2e tests (requires cluster)
make build              # build binary
make image              # build docker image
```

## Contributing

1. Fork the repo and create a feature branch.
2. Make your changes and run `make verify` to ensure everything passes.
3. Open a pull request with a clear description of the change.

For significant changes, please open an issue first to discuss the approach.

## License

[Apache License 2.0](LICENSE)
