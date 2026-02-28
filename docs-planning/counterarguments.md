# Counterarguments: Why NOT to Build a Documentation Site Right Now

## Executive Summary

Kelos is at v0.14.0 with 2 active contributors, ~48 GitHub stars, and a rapidly evolving API surface. Building a dedicated documentation site right now is **premature optimization** that will create a maintenance burden disproportionate to the project's current scale and contributor capacity. The existing documentation is already surprisingly comprehensive. There are higher-ROI investments that would serve users better.

---

## 1. Premature Optimization — The Numbers Don't Justify It

### Project Scale (as of Feb 2026)

| Metric | Value |
|--------|-------|
| Version | v0.14.0 (v1alpha1 API) |
| GitHub stars | ~48 |
| Git contributors | 2 (Gunju Kim, Aslak Knutsen) |
| Total commits | 450 |
| Go source files | 130 |
| CRD type definition lines | 745 (across 4 `*_types.go` files) |

### The existing docs are already substantial

| File | Lines | Content |
|------|-------|---------|
| `README.md` | 583 | Quick start, how it works, examples, security, cost, FAQ |
| `docs/reference.md` | 310 | Full CRD field reference, CLI reference, config reference |
| `docs/agent-image-interface.md` | 132 | Custom image interface spec |
| Example READMEs (7 dirs) | ~400+ | Step-by-step walkthroughs for each pattern |
| YAML examples | 25 files | Ready-to-apply manifests |
| **Total** | **~1,400+ lines** | Covers nearly every feature |

**The README alone is a near-complete user guide.** It includes: Quick Start (with `kind` setup), How It Works (with architecture diagram), 5 runnable example patterns with YAML, orchestration patterns, security considerations, cost management, FAQ, CLI reference, and development guide. A doc site would mostly be re-organizing this same content with navigation chrome.

### Cost/benefit at this scale

At ~48 stars, the user base is likely under 20 active users. The question isn't "would a doc site be nice?" — it's "would a doc site be a better use of 2 contributors' time than anything else?" Every hour spent on a doc site is an hour NOT spent on:
- Shipping features that attract users
- Fixing bugs that retain users
- Writing examples that teach users
- Improving error messages that unblock users

**Doc sites become ROI-positive when the project has enough users that the docs scale better than individual support.** At ~48 stars, personal responses in GitHub issues likely still outperform any documentation.

---

## 2. Maintenance Cost — The Surface Area Is Already Large

### CRD Field Count (from `api/v1alpha1/*_types.go`)

I counted every struct field that would need documenting:

| Struct | Fields |
|--------|--------|
| `TaskSpec` | 11 (type, prompt, credentials, model, image, workspaceRef, agentConfigRef, dependsOn, branch, ttlSecondsAfterFinished, podOverrides) |
| `TaskStatus` | 8 (phase, jobName, podName, startTime, completionTime, message, outputs, results) |
| `PodOverrides` | 4 (resources, activeDeadlineSeconds, env, nodeSelector) |
| `Credentials` | 2 (type, secretRef) |
| `WorkspaceSpec` | 5 (repo, ref, secretRef, remotes, files) |
| `AgentConfigSpec` | 3 (agentsMD, plugins, mcpServers) |
| `PluginSpec` | 3 (name, skills, agents) |
| `MCPServerSpec` | 7 (name, type, command, args, url, headers, env) |
| `TaskSpawnerSpec` | 6 (when, taskTemplate, pollInterval, maxConcurrency, suspend, maxTotalTasks) |
| `GitHubIssues` | 10 (repo, types, labels, excludeLabels, state, triggerComment, excludeComments, assignee, author, priorityLabels) |
| `Jira` | 4 (baseUrl, project, jql, secretRef) |
| `TaskTemplate` | 11 (type, credentials, model, image, workspaceRef, agentConfigRef, dependsOn, branch, promptTemplate, ttlSecondsAfterFinished, podOverrides) |
| `TaskSpawnerStatus` | 9 (phase, deploymentName, cronJobName, totalDiscovered, totalTasksCreated, activeTasks, lastDiscoveryTime, message, conditions) |
| **Total** | **~83 fields** |

### Concepts requiring documentation

Beyond field references, a doc site would need conceptual pages for:

1. **4 CRD resources** — Task, Workspace, AgentConfig, TaskSpawner
2. **4 agent types** — claude-code, codex, gemini, opencode (each with distinct credential flows)
3. **3 event source types** — GitHub Issues, Cron, Jira
4. **2 authentication methods** — PAT vs GitHub App (with secret structure differences)
5. **2 credential types** — API key vs OAuth
6. **Dependency chaining** — `dependsOn`, result passing via Go templates, `.Deps` map structure
7. **Branch serialization** — mutex behavior, interaction with `dependsOn`
8. **MCP servers** — 3 transport types (stdio, http, sse), each with different config fields
9. **Plugins** — skills, agents, plugin directories
10. **Agent Image Interface** — entrypoint, env vars, UID, output capture protocol
11. **CLI** — ~10 commands with multiple flags each
12. **Config file** — `~/.kelos/config.yaml` with its own schema
13. **Orchestration patterns** — self-development, fork workflows, fleet refactoring
14. **Security model** — pod isolation, token scoping, RBAC
15. **Cost management** — model selection, maxConcurrency, timeouts

**That's 15+ conceptual topics and 83+ reference fields.** With the API still at `v1alpha1`, these fields WILL change. Every change requires updating both the code AND the docs. Two contributors cannot sustain this.

### The v1alpha1 problem

The API is explicitly `v1alpha1` — it signals instability. Recent commits show the API is still actively evolving:
- `628dfd5` — Fixed incorrect field names in example workspace.yaml (docs were already wrong)
- `0572e18` — Added GitHub App authentication (new Workspace fields)
- `83d2dae` — Changed cron-based TaskSpawners to use CronJob (new status fields)
- `6052977` — Added TriggerComment support (new GitHubIssues fields)

A doc site amplifies the cost of every API change. Without one, you update `README.md` and `docs/reference.md`. With a doc site, you update the source files, the site content, verify cross-links, check rendering, and deploy.

---

## 3. Better Alternatives — Higher ROI Investments

The same engineering effort could go into things that help users more:

### A. Better README organization (2-4 hours)

The README is 583 lines and covers everything, but it's getting long. A lightweight restructuring — splitting the Quick Start into its own file, moving the CLI reference to `docs/` — would improve discoverability without the overhead of a site.

### B. More examples (4-8 hours per example)

The 7 examples are the project's best documentation. They're concrete, runnable, and self-contained. Adding examples for common pain points (Jira integration, multi-repo fan-out, GitHub Enterprise setup, using custom images) would directly help users more than a doc site.

### C. Better error messages (ongoing)

When `628dfd5` fixed incorrect field names in an example, the real question is: **why didn't the controller's error message tell the user what was wrong?** Investing in validation webhooks and clear error messages eliminates an entire class of documentation needs. The best documentation is an error message that tells you exactly what to fix.

### D. Inline godoc (2-4 hours)

The Go types already have decent comments (e.g., `// Branch is the git branch this Task works on. The controller ensures only one Task with the same Branch value runs at a time.`). Improving these and publishing via pkg.go.dev costs nearly zero maintenance — godoc stays in sync with the code by definition.

### E. Troubleshooting guide (2-4 hours)

A single `docs/troubleshooting.md` addressing common issues (credential errors, workspace not found, rate limits, pod failures) would provide more value per line than any doc site page.

---

## 4. Risks of Building Too Early

### Stale documentation

Stale docs are worse than no docs. Users who follow outdated instructions waste time and lose trust. With `v1alpha1` APIs changing every few releases, the probability of staleness is high.

Evidence this is already happening:
- Commit `628dfd5` (the most recent non-merge commit): "Fix incorrect field names in example 05 workspace.yaml" — the examples themselves, which live right next to the code, already drifted. A doc site living in a separate repo or build pipeline would drift faster.

### Version drift

When docs and code diverge, users can't tell which to trust. A doc site creates a second "source of truth" that can silently fall behind. The README/reference approach keeps docs close to the code where they're more likely to be updated during the same PR.

### Split attention

Two contributors means every task has high opportunity cost. Time spent on docs infrastructure (choosing a framework, designing the site, setting up CI/CD for docs, configuring a domain) is time not spent on the product itself. For a pre-1.0, low-adoption project, product improvement is almost always more impactful than documentation polish.

### Premature terminology commitment

A doc site forces you to commit to terminology and conceptual models. If you write a "Concepts" page explaining that Tasks are "ephemeral units of work," but later realize Tasks should be persistent, you've created an artifact that resists change. Right now, the terminology is in the code and can evolve freely.

### The "empty restaurant" effect

A doc site with sparse content looks worse than no doc site at all. If you launch with 5 pages where 3 are stubs or "TODO" sections, users perceive the project as unfinished or abandoned. The current README gives an impression of completeness because it IS complete for the current scope.

---

## 5. Minimum Viable Improvement — What to Do Instead

If the goal is "make Kelos easier to learn," here's the smallest investment that meaningfully helps users WITHOUT a full doc site:

### Tier 1: Zero-maintenance improvements (1-2 days)

1. **Split the README.** Move "Cost and Limits" and "Security Considerations" to `docs/`. Keep the README focused on Quick Start, Why Kelos, and Examples. Add a clean "Documentation" section with links.
2. **Add a `docs/troubleshooting.md`.** Cover the top 5 failure modes users hit (bad credentials, missing workspace, rate limits, pod OOM, branch conflicts).
3. **Improve CLI `--help` text.** Make `kelos run --help` comprehensive enough that users rarely need the docs. This is zero-maintenance because it ships with the binary.

### Tier 2: Low-maintenance improvements (3-5 days)

4. **Add 2-3 more examples.** Jira integration, GitHub Enterprise, and custom agent images are the obvious gaps.
5. **Expand `docs/reference.md`** with a "Common Patterns" section showing real-world YAML snippets (not just field tables).
6. **Add validation webhooks** that produce clear error messages for the most common misconfigurations.

### Tier 3: Only if user demand proves it (1-2 weeks)

7. **A simple doc site** — but ONLY after reaching ~200+ stars, 5+ contributors, or receiving repeated requests for better docs in GitHub issues. Use the simplest possible setup (GitHub Pages + a single-page generator) and migrate the existing markdown files with minimal changes.

---

## Conclusion

Building a documentation site for Kelos right now is solving a problem the project doesn't yet have. The existing README + reference docs + 7 examples already cover nearly the entire feature surface. The two active contributors' time is better spent on product improvements, better error messages, and more examples. A doc site should be triggered by user demand (GitHub issues asking "where are the docs?"), not by aspirational project maturity.

**The first documentation site is free. Keeping it accurate costs forever.**
