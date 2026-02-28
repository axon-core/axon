# Documentation Audit: Kelos Project

## Files Reviewed

- `README.md` (584 lines)
- `docs/reference.md` (310 lines)
- `docs/agent-image-interface.md` (133 lines)
- `examples/README.md` and all 7 example READMEs + YAML files
- `self-development/README.md` (272 lines) + key YAML files
- `CLAUDE.md` / `AGENTS.md` (identical, 24 lines each)
- `api/v1alpha1/task_types.go`, `taskspawner_types.go`, `workspace_types.go`, `agentconfig_types.go`

---

## 1. What's Working Well

### Strong README Structure
The `README.md` is well-organized with a clear hierarchy: Demo → Why Kelos → Quick Start → How It Works → Examples → Reference → Security → Cost/Limits → FAQ. This follows the "inverted pyramid" pattern well — users get the value proposition immediately.

### Excellent Quick Start
The Quick Start section (README.md:54-213) gets users from zero to a running task in 4 steps. The `kind` cluster hint in a collapsible section is a nice touch. The CLI-first approach (`kelos run -p "..."`) is much more accessible than starting with raw YAML.

### Examples Are Genuinely Useful
The 7 examples form a logical progression:
1. Simple task (no workspace)
2. Task with workspace (git repo)
3. TaskSpawner for issues
4. TaskSpawner with cron
5. AgentConfig injection
6. Fork workflow
7. Task pipeline

Each example README follows a consistent template: Use Case → Resources table → Steps → Notes. The YAML files include `# TODO:` comments for placeholders — this is a good pattern for copy-paste-edit workflows.

### Good "Why" Section
README.md:44-52 clearly articulates five differentiators. Each bullet leads with a bolded phrase and expands with specifics. The comparison framing ("Orchestration, not just execution") is effective.

### Self-Development as Showcase
The `self-development/` directory (README.md:415-425, self-development/README.md) is a compelling real-world example. It demonstrates the product's own capabilities by showing how Kelos uses itself for autonomous development. The feedback loop pattern (`excludeLabels: [kelos/needs-input]`) is well-explained.

### Security Section is Responsible
README.md:463-481 takes a responsible approach — it's clear about what agents CAN and CANNOT do, and the `--dangerously-skip-permissions` explanation (line 479) proactively addresses a common concern.

---

## 2. What's Confusing, Buried, or Duplicated

### Reference Doc is Incomplete — Significant Gaps vs. Code

The `docs/reference.md` is the authoritative API reference, but it's missing several fields that exist in the Go types:

**Workspace — missing fields:**
- `spec.remotes[]` (GitRemote) — defined in `workspace_types.go:8-17`, used in example 06-fork-workflow, but completely absent from `docs/reference.md`. The Workspace table only lists `repo`, `ref`, `secretRef.name`, and `files[]`.

**TaskSpawner GitHubIssues — missing fields:**
- `triggerComment` (taskspawner_types.go:80-87) — comment-based discovery trigger
- `excludeComments` (taskspawner_types.go:89-95) — comment-based exclusion
- `assignee` (taskspawner_types.go:97-101) — filter by assignee
- `author` (taskspawner_types.go:103-107) — filter by issue creator
- `priorityLabels` (taskspawner_types.go:109-115) — priority ordering for task creation
- `repo` override (taskspawner_types.go:51-57) — override which repo to poll

None of these six fields appear anywhere in `docs/reference.md`. The `priorityLabels` field is actively used in `self-development/kelos-workers.yaml:13-17` but never documented.

**TaskSpawner Jira source — entirely undocumented:**
- The `Jira` struct (taskspawner_types.go:118-146) defines a complete Jira integration with `baseUrl`, `project`, `jql`, and `secretRef`. This is not mentioned anywhere in any documentation file — not in the README, not in the reference, not in any example.

### Credential Configuration is Scattered and Confusing

Credential setup is spread across at least 5 locations:
1. README.md Quick Start (lines 111-153) — config.yaml format
2. README.md YAML section (lines 174-226) — Kubernetes Secret + Task YAML
3. docs/reference.md Configuration section (lines 178-198) — config.yaml fields
4. docs/reference.md Workspace Authentication (lines 53-87) — workspace secrets
5. docs/agent-image-interface.md (lines 27-46) — environment variable mapping

A newcomer trying to understand "how do I set up credentials?" would need to read all five locations and mentally merge them. The distinction between agent credentials (Anthropic/OpenAI key) and workspace credentials (GitHub token) is never explicitly called out as two separate concepts.

### `token` vs `secretRef` vs `oauthToken` Naming Confusion

The config.yaml uses `oauthToken` and `token` (README.md:114-118):
```yaml
oauthToken: <your-oauth-token>
workspace:
  token: <github-token>
```

The YAML resources use `secretRef` and `credentials.secretRef`:
```yaml
credentials:
  secretRef:
    name: claude-oauth-token
workspace:
  secretRef:
    name: github-token
```

These are the same concepts but with different names depending on whether you use CLI/config or YAML. This mapping is never explicitly documented.

### AgentConfig `plugins` vs Workspace `files` Overlap

Both AgentConfig (agentconfig_types.go:15-19) and Workspace (workspace_types.go:55-60) can inject files into the agent's environment:
- AgentConfig `plugins` → mounted as plugin directories
- Workspace `files[]` → written into the cloned repo

The docs mention both but never explain when to use which. The Workspace `files` description says it can inject "plugin-like assets such as skills" — this directly overlaps with AgentConfig's purpose. A newcomer would not know which mechanism to choose.

### "How It Works" Diagram is an Image, Not Text

README.md:232 references an image hosted on GitHub user-attachments. This means:
- It's not searchable
- It can't be updated without re-uploading
- It's not accessible (no alt text beyond "kelos-resources")

### Duplicate Template Variable Tables

The promptTemplate variables table appears in THREE places:
1. `docs/reference.md:136-149`
2. `self-development/README.md:210-221`
3. README.md doesn't have the table but describes `.Deps` at line 361-382

The tables in reference.md and self-development/README.md are copies. If one is updated, the other risks going stale.

---

## 3. Gaps Between Code and Docs

### Jira Integration — Zero Documentation
The biggest gap. `taskspawner_types.go:118-146` defines a full Jira source with:
- `baseUrl`, `project`, `jql` fields
- Secret-based auth (JIRA_TOKEN, optional JIRA_USER)
- Support for both Jira Cloud (Basic auth) and Data Center (Bearer token)

There is no mention of Jira anywhere in the README, reference, examples, or FAQ. A user reading the docs would have no idea Jira is supported.

### Fork Workflow — Advanced but Under-documented
Example 06 demonstrates the fork workflow, but the key enabler — `spec.when.githubIssues.repo` override — is not documented in `docs/reference.md`. Similarly, `workspace.spec.remotes` is used in the example but absent from the reference.

### Comment-Based Triggers — Powerful but Hidden
`triggerComment` and `excludeComments` (taskspawner_types.go:79-95) enable a comment-based workflow for repos where you lack label permissions. This is a significant feature for open-source contribution scenarios but is completely undocumented.

### GitHub Enterprise Support — Mentioned Only in Env Vars
`docs/agent-image-interface.md:41` mentions `GH_ENTERPRISE_TOKEN` and `GH_HOST` environment variables for GitHub Enterprise, but this capability is never mentioned in the main README or reference doc. An enterprise user would not discover this.

### Task Immutability
`task_types.go:183` has a validation rule: `"Task spec is immutable after creation"`. This is never documented. A user trying to update a running task's prompt would get an opaque error.

### `kelos create workspace` and `kelos create agentconfig` CLI Commands
These appear in the CLI reference table (reference.md:273-274) but have no examples, no flags documented, and no explanation. A user seeing them would have to guess at usage.

---

## 4. Discoverability for Newcomers

### Good: Top-Down Entry Point
A newcomer landing on the README gets a clear path: Demo → Why → Quick Start → Examples. The Quick Start's 4-step flow is accessible.

### Problem: No Conceptual Overview
There is no document that explains the mental model. The "Core Primitives" section (README.md:238-258) lists the four resources but doesn't explain how they relate or when you'd use each one. Questions like:
- "When do I need a Workspace vs. just running without one?"
- "What's the difference between AgentConfig and putting instructions in the prompt?"
- "When should I use TaskSpawner vs. just creating Tasks?"

These require reading multiple sections and inferring the answers.

### Problem: No Troubleshooting Guide
The self-development README has a "Troubleshooting" section (self-development/README.md:249-264), but it's buried inside self-development, which most newcomers won't explore. The main README has no troubleshooting section — only a single tip about controller logs (README.md:171-172).

### Problem: No "What Happens When Things Go Wrong"
If a task fails, the docs never explain:
- How to view the failure reason (`status.message`, `status.phase`)
- How to retry a failed task (delete and recreate? Is there a retry mechanism?)
- What the Pod logs contain and how to interpret them
- Common failure modes (rate limiting, invalid credentials, Git auth failures)

### Problem: Agent Type Differences Aren't Compared
The docs mention four agent types (claude-code, codex, gemini, opencode) but never compare them. A newcomer choosing between agents has no guidance on:
- Which features work with which agents (plugins only work with claude-code)
- Credential differences between agents
- Capability differences
- Cost/performance tradeoffs

---

## 5. Progression from Simple to Complex

### Good: Examples Follow a Learning Curve
The 7 examples progress logically:
- 01: Minimal (just a task + secret)
- 02: Add workspace (git repo)
- 03: Add automation (TaskSpawner)
- 04: Add scheduling (Cron)
- 05: Add configuration (AgentConfig)
- 06: Advanced pattern (Fork workflow)
- 07: Advanced pattern (Pipeline with dependencies)

### Problem: No Intermediate "How-To" Layer
There's a gap between "here's a working example" and "here's the full API reference." Missing middle layer would include:
- How to set up cost controls (maxConcurrency, timeouts)
- How to debug a failing agent
- How to structure prompts for best results
- How to set up branch protection with Kelos
- How to migrate from CLI-only to YAML-managed

### Problem: Self-Development is Aspirational, Not Instructional
The `self-development/` directory is impressive but serves more as a showcase than a learning resource. The prompts in `kelos-workers.yaml` (83 lines of promptTemplate) are highly Kelos-specific. A newcomer wanting to build their own autonomous workflow doesn't get intermediate steps between "simple task" and "full self-development pipeline."

### Problem: Reference Doc Assumes Familiarity
`docs/reference.md` is a flat list of fields. It works well as a lookup table for someone who already understands the system, but it doesn't help someone learning. There are no:
- Explanatory notes about common patterns
- Cross-references between related fields
- Warnings about common mistakes
- Examples within the reference itself

---

## Summary: Overall Documentation Health

| Aspect | Rating | Notes |
|--------|--------|-------|
| Getting started | Good | Quick Start is accessible; CLI-first approach works well |
| API completeness | Poor | At least 8 fields undocumented; Jira integration entirely missing |
| Example quality | Good | Consistent format, logical progression, runnable YAML |
| Conceptual clarity | Fair | Mental model must be inferred; no explicit concept docs |
| Troubleshooting | Poor | Essentially absent from main docs |
| Discoverability | Fair | Good entry point, but advanced features are hidden |
| Maintenance burden | Concerning | Triple-duplicated tables, image-only diagrams, scattered credential docs |
