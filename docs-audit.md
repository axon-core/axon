# Documentation Audit: Kelos Project

## Executive Summary

Kelos has **strong documentation for an early-stage project** — the README is well-structured, examples are thorough, and the reference doc is comprehensive. However, the documentation has significant **discoverability, fragmentation, and gap problems** that will worsen as the project grows. Key findings:

- **2 doc files** (`docs/reference.md`, `docs/agent-image-interface.md`) carry the entire reference burden
- **3 major API features** are undocumented or barely mentioned (Jira integration, comment-based triggers, priority labels)
- **Information is scattered** across README, docs/, examples/, and self-development/ with no navigation or search
- **No conceptual guides** exist — only reference tables and examples

---

## 1. What's Working Well

### 1.1 README.md is excellent for its purpose
The README (`README.md`, 584 lines) serves as both landing page and getting-started guide. Strengths:

- **Clear value proposition** (lines 1-20): One-liner + supported agents + demo in the first screenful
- **Quick Start is genuinely quick** (lines 54-213): 4 numbered steps, collapsible alternatives, copy-pasteable commands
- **Real-world examples inline** (lines 261-425): TaskSpawner, dependency chains, AgentConfig, self-development pipeline — all with both CLI and YAML equivalents
- **Security and cost sections** (lines 463-519): Honest, practical advice with concrete YAML snippets for `maxConcurrency`, `activeDeadlineSeconds`, `suspend`
- **Progressive disclosure via `<details>` blocks**: kubectl alternative (line 174), API key auth (line 215), CLI reference (line 444) — keeps the main flow clean

### 1.2 Examples are well-structured and progressive
The `examples/` directory (7 examples) follows a clear learning progression:

| Example | Builds On | New Concept |
|---------|-----------|-------------|
| 01-simple-task | — | Minimal Task + API key |
| 02-task-with-workspace | 01 | Git workspace + OAuth |
| 03-taskspawner-github-issues | 02 | Auto-spawn from GitHub issues |
| 04-taskspawner-cron | 02 | Cron-triggered tasks |
| 05-task-with-agentconfig | 02 | Reusable config + plugins |
| 06-fork-workflow | 03 | Fork/upstream pattern |
| 07-task-pipeline | 02 | dependsOn + result passing |

Each example has:
- A README with use case, resource table, step-by-step instructions, and cleanup
- All necessary YAML files with `# TODO:` placeholders
- `examples/README.md` ties them together with a summary table

### 1.3 Reference doc covers the API surface well
`docs/reference.md` (310 lines) provides field-level documentation for all 4 CRDs plus CLI flags. It accurately reflects the Go types with Required/Optional annotations and descriptions.

### 1.4 Agent Image Interface is thorough
`docs/agent-image-interface.md` (133 lines) is a solid technical spec for custom image authors. The environment variable table (lines 29-46), output capture protocol (lines 68-101), and entrypoint pattern (lines 109-125) give everything needed to build a compatible image.

### 1.5 Self-development README is a strong showcase
`self-development/README.md` (272 lines) documents 5 TaskSpawner patterns with clear tables, deploy commands, and a well-explained "feedback loop" pattern. It doubles as both documentation and proof that Kelos works in production.

---

## 2. What's Confusing, Buried, Duplicated, or Missing

### 2.1 Duplicated content across files

**promptTemplate variables table** appears in 3 places with identical content:
- `docs/reference.md` lines 136-149
- `self-development/README.md` lines 210-221
- `README.md` lines 361-382 (partial, in prose form)

If any of these get updated, the others will drift.

**Workspace authentication** is explained in:
- `README.md` lines 121-153 (inline in Quick Start)
- `docs/reference.md` lines 53-87 (formal reference)
- `examples/02-task-with-workspace/README.md` lines 54-64 (GitHub App note)
- `self-development/README.md` lines 16-37 (prerequisite)

Four slightly different explanations of the same concept.

**CLI commands** appear in:
- `README.md` lines 444-461 (collapsed summary table)
- `docs/reference.md` lines 254-309 (full reference with flags)

### 2.2 Missing conceptual documentation

There are **no conceptual guides** explaining:

- **Architecture overview**: How does the controller work? What's the reconciliation loop? What Kubernetes resources does it create under the hood (Jobs, Deployments, CronJobs)?
- **Task lifecycle**: The phases (Pending → Waiting → Running → Succeeded/Failed) are mentioned in the reference table (`docs/reference.md` line 155) but never explained with a diagram or state machine.
- **Branch serialization**: Mentioned in `README.md` line 321 ("Tasks sharing the same branch are serialized automatically") and `examples/07-task-pipeline/README.md` lines 62-64, but never fully explained. What happens when a second task targets the same branch? Does it queue? Does it fail?
- **How credentials flow**: The agent image interface doc lists env vars, but there's no end-to-end explanation of how a secret becomes an env var in the agent pod.
- **Troubleshooting guide**: Only `self-development/README.md` lines 249-264 has troubleshooting content, and it's specific to self-development. No general troubleshooting guide exists.

### 2.3 Buried information

- **GitHub Enterprise support**: Only discoverable by reading `docs/agent-image-interface.md` lines 40-41 (`GH_ENTERPRISE_TOKEN`, `GH_HOST`). Not mentioned in README or reference.
- **`spec.files[]` on Workspace**: Listed in `docs/reference.md` line 51 and defined in `api/v1alpha1/workspace_types.go` lines 19-29, but never shown in any example. The Go type comment explains it can inject skills and instruction files — a useful feature with zero documentation beyond the field table.
- **Task immutability**: The Go type has a CEL validation rule (`api/v1alpha1/task_types.go` line 183: `rule="self == oldSelf",message="Task spec is immutable after creation"`) — important operational knowledge that's not mentioned anywhere in docs.
- **The `kelos create workspace` and `kelos create agentconfig` commands**: Listed in `docs/reference.md` line 273-274 but never demonstrated in any example or tutorial.
- **`workspace.name` vs `workspace.repo` in config**: `docs/reference.md` lines 203-243 explains the two forms, but the Quick Start only shows the inline `repo` form. Users might not discover they can reference an existing Workspace by name.

### 2.4 Confusing elements

- **CLAUDE.md vs AGENTS.md**: Both exist at the repo root. `CLAUDE.md` and `AGENTS.md` have **identical content** (same 24 lines). It's unclear why both exist or which takes precedence. The AgentConfig docs reference writing to `~/.claude/CLAUDE.md`, adding to the confusion.
- **`oauthToken: "@~/.codex/auth.json"` syntax** (`README.md` line 133): The `@` prefix for file references is never explained. Is this a Kelos convention? A Codex convention? What other fields support it?
- **"Type" field naming collision**: `spec.type` on Task means "agent type" (claude-code, codex, etc.), while `spec.credentials.type` means credential type (api-key, oauth), and `spec.mcpServers[].type` means transport type (stdio, http, sse). Same field name, three different meanings.
- **Config file vs YAML resources**: The Quick Start uses `~/.kelos/config.yaml` with `oauthToken` and `workspace.token`, while the YAML examples use Kubernetes Secrets. The relationship between these two approaches is never clearly explained. When does the CLI auto-create secrets vs when must you create them manually?

---

## 3. Gaps Between Code Capabilities and Documentation

### 3.1 Jira integration — COMPLETELY UNDOCUMENTED

The `api/v1alpha1/taskspawner_types.go` lines 118-146 define a full `Jira` struct:

```go
type Jira struct {
    BaseURL   string          `json:"baseUrl"`
    Project   string          `json:"project"`
    JQL       string          `json:"jql,omitempty"`
    SecretRef SecretReference `json:"secretRef"`
}
```

With detailed comments about Cloud vs Data Center auth patterns. This feature is:
- Not mentioned in README.md
- Not mentioned in docs/reference.md
- Has no example in examples/
- Has no mention anywhere in the entire documentation

### 3.2 Comment-based triggers — UNDOCUMENTED

`api/v1alpha1/taskspawner_types.go` lines 79-95 define `TriggerComment` and `ExcludeComments`:

```go
TriggerComment  string   `json:"triggerComment,omitempty"`
ExcludeComments []string `json:"excludeComments,omitempty"`
```

These enable a slash-command pattern (e.g., `/kelos pick-up`) for repos where you lack label permissions. The Go comments explain the interaction between trigger and exclude comments well, but this is **nowhere in the documentation**.

### 3.3 Priority labels — UNDOCUMENTED

`api/v1alpha1/taskspawner_types.go` lines 109-115 define `PriorityLabels`:

```go
PriorityLabels []string `json:"priorityLabels,omitempty"`
```

This enables label-based prioritization when `maxConcurrency` limits task creation. Not documented anywhere.

### 3.4 Assignee and Author filters — UNDOCUMENTED

`api/v1alpha1/taskspawner_types.go` lines 97-107 define `Assignee` and `Author` fields with server-side GitHub API filtering. Not in `docs/reference.md` TaskSpawner table.

### 3.5 Workspace `Remotes` field — UNDOCUMENTED in reference

`api/v1alpha1/workspace_types.go` lines 8-17 define `GitRemote` and lines 48-53 define `Remotes` on WorkspaceSpec with CEL validation (no "origin" name, unique names). The field:
- Is used in `examples/06-fork-workflow/workspace.yaml`
- Is explained in `examples/06-fork-workflow/README.md`
- Is **NOT** listed in `docs/reference.md` Workspace table

### 3.6 GitHub Issues `Types` field — partially documented

The `Types` field (`api/v1alpha1/taskspawner_types.go` lines 59-63) allows filtering for `"issues"`, `"pulls"`, or both. It appears in `docs/reference.md` line 115 but has no example showing how to use it for PR-triggered workflows.

### 3.7 Workspace `Files` field — in reference but no example

Listed in `docs/reference.md` line 51 (`spec.files[]`), defined in `api/v1alpha1/workspace_types.go` lines 19-29. The Go comment explains use cases (injecting skills, CLAUDE.md) but no example demonstrates it. This is distinct from AgentConfig's plugins — `files` is a workspace-level injection mechanism.

---

## 4. Discoverability Assessment

### 4.1 No search capability

With only 2 files in `docs/`, finding information requires knowing which file to look in. A user wondering "how do I set up Jira integration" would find nothing. A user searching for "timeout" would need to check README, reference, and examples separately.

### 4.2 No navigation or table of contents

- `docs/` has no index file
- README.md links to `docs/reference.md` and `docs/agent-image-interface.md` but there's no way to discover all available docs
- `self-development/README.md` is a rich source of patterns but is not linked from `docs/` at all
- No sitemap, no sidebar, no breadcrumbs

### 4.3 No clear reading path for different audiences

Three distinct user types exist:
1. **New user**: Wants Quick Start → first task → first workspace → first TaskSpawner
2. **Operator**: Wants security, cost controls, RBAC, monitoring, troubleshooting
3. **Custom agent builder**: Wants the agent image interface spec

Currently all three audiences must navigate the same README and hope they find the relevant `<details>` block.

### 4.4 Cross-reference quality

Good cross-references:
- README → `docs/reference.md` sections (lines 437-442)
- README → `docs/agent-image-interface.md` (lines 22, 50)
- README → `examples/` directory (line 425)
- Example READMEs → `docs/reference.md` sections

Missing cross-references:
- No link from anywhere to `self-development/` except README line 421
- `docs/reference.md` never links to relevant examples
- `docs/agent-image-interface.md` never links to the entrypoint source files it references (line 129-132)

### 4.5 Information density

The README is 584 lines. While well-organized with collapsible sections, it carries too much weight — Quick Start, concepts, examples, security, cost, FAQ, CLI reference, and development all in one file. Users scrolling past the Quick Start to find "how do I chain tasks" must scan ~200 lines of content.

---

## 5. Summary of Documentation Inventory

| Location | Files | Lines | Purpose |
|----------|-------|-------|---------|
| `README.md` | 1 | 584 | Landing page + Quick Start + examples + security + FAQ |
| `docs/` | 2 | 443 | Reference tables + agent image spec |
| `examples/` | 7 dirs | ~380 (READMEs) + ~250 (YAML) | Progressive examples with step-by-step guides |
| `self-development/` | 1 README + 6 YAML | 272 + ~600 | Real-world orchestration patterns |
| `CLAUDE.md` / `AGENTS.md` | 2 | 24 each | AI assistant conventions (identical content) |
| **Total** | ~20 files | ~2,600 lines | — |

### Feature documentation coverage

| Feature | Documented? | Where |
|---------|------------|-------|
| Task (basic) | Yes | README, reference, examples 01-02 |
| Workspace (basic) | Yes | README, reference, examples 02 |
| TaskSpawner (GitHub Issues) | Yes | README, reference, examples 03 |
| TaskSpawner (Cron) | Yes | README, reference, examples 04 |
| AgentConfig | Yes | README, reference, examples 05 |
| dependsOn / pipelines | Yes | README, reference, examples 07 |
| Fork workflow | Yes | Examples 06 |
| Custom agent images | Yes | docs/agent-image-interface.md |
| Self-development patterns | Yes | self-development/README.md |
| CLI reference | Yes | README, reference |
| Workspace `remotes` | Partial | Example 06 only, missing from reference |
| Workspace `files` | Partial | Reference table only, no example |
| GitHub Enterprise | Barely | Agent image interface env var table only |
| Jira integration | **No** | Only in Go types |
| TriggerComment / ExcludeComments | **No** | Only in Go types |
| PriorityLabels | **No** | Only in Go types |
| Assignee / Author filters | **No** | Only in Go types |
| Task immutability | **No** | Only in Go CEL validation rule |
| Config file ↔ YAML resource relationship | **No** | Fragments across README + reference |
