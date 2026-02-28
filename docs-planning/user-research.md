# User Perspective Research: Kelos Documentation

## 1. User Personas

Based on analysis of README.md, examples/, self-development/, and GitHub issues.

### Persona A: Platform Engineer / Internal Developer Platform Builder (Primary)

**Evidence:**
- README emphasizes scalable parallelism, fan-out across repos, Kubernetes scheduling/resource management (line 51)
- GitHub App auth recommended "for organizations" and "production/org use" (lines 142-152)
- `maxConcurrency`, `maxTotalTasks`, namespace RBAC isolation, scoped ServiceAccounts — organizational-scale controls
- Fleet-wide refactoring and AI Worker Pools explicitly named as use cases (lines 431-433)
- ArgoCD/GitOps integration highlighted

**K8s experience:** High. Comfortable with CRDs, RBAC, ServiceAccounts, namespace isolation, ArgoCD.

**Workflow:** Deploy Kelos into shared cluster, configure GitHub App auth for the org, define reusable Workspaces and AgentConfigs, let developers trigger tasks via `kelos run` or declarative YAML.

### Persona B: DevOps / Automation Engineer

**Evidence:**
- "Hands-Free CI/CD" use case (line 432)
- "Event-Driven Bug Fixing" via GitHub issue labels (line 430)
- Cron-triggered TaskSpawner (example 04)
- self-development pipeline with `kelos-triage.yaml` at concurrency 8

**K8s experience:** Intermediate. Uses existing cluster infrastructure but may not build clusters from scratch.

**Workflow:** Wire GitHub issue labels to TaskSpawners, configure cron schedules, integrate outputs into Slack/CI dashboards.

### Persona C: Individual Developer / Indie Hacker (Aspirational)

**Evidence:**
- Quick Start framing: "Get running in 5 minutes"
- `kelos run -p "Fix the bug..."` one-liner demo (line 276)
- `kind create cluster` as local option (line 68)
- FAQ explicitly addresses "Can I use Kelos without Kubernetes?" (line 530)

**K8s experience:** Low to intermediate.

**Friction:** Highest. They want `pip install && run` but get `kind create cluster && kelos install && kubectl apply -f secret.yaml`. The Kubernetes requirement is a hard gate.

### Persona D: Open-Source Maintainer

**Evidence:**
- Example 06 "Fork Workflow" — for contributors without upstream push access
- `kelos-workers.yaml` includes "Check if PR already exists" logic
- self-development/ shows Kelos managing its own contributions

**K8s experience:** Intermediate.

---

## 2. User Journey: Discovery to First Task

### Step-by-step path and friction points

| Step | Action | Friction Level | Notes |
|------|--------|---------------|-------|
| 0. Discovery | Land on README | Low | Clear tagline, demo section, video embed |
| 1. Prerequisites | Need K8s cluster (1.28+) | **HIGH** | Hard-blocked without K8s. `kind` escape hatch exists but adds significant overhead |
| 2. Install CLI | `curl ... \| bash` or `go install` | Low | Single binary, straightforward |
| 3. Install to cluster | `kelos install` | Medium | Requires kubectl + cluster admin permissions |
| 4. Configure | `kelos init` + edit config | **HIGH** | 4 credential types to understand (Claude OAuth, API key, Codex OAuth, GitHub PAT/App). No "start here" recommendation |
| 5. Run first task | `kelos run -p "..."` | Medium | Works, but ephemeral warning appears *after* setup. Results disappear without workspace config |
| 6. See results | `kelos logs` | Medium | No next-step guidance after task creation (issue #183) |

### Critical friction points

1. **Cluster requirement gate** — FAQ acknowledges users will ask "Can I use this without K8s?" Answer is no. `kind` works locally but adds significant setup.

2. **Credential complexity** — Before the first task, users must understand: AI agent credentials (OAuth vs API key), GitHub credentials (PAT vs GitHub App), and workspace config. The README presents these in parallel without "start here" guidance.

3. **Step ordering contradiction** — Issue #305: The CLI's `kelos init` suggests order `init -> install -> run`, but README shows `install -> init -> run`. Directly contradicts what users see.

4. **OAuth vs API key confusion** — Issue #201: `--credential-type` defaults to `api-key` in CLI code, but Quick Start presents OAuth first. Users picking one type may get silent failures.

5. **Ephemeral result surprise** — Warning "Without a workspace, the agent runs in an ephemeral pod — any files it creates are lost" appears after config setup, not before. Users running example 01 (no workspace) will wonder why there's no git output.

6. **No debugging path** — Main troubleshooting guidance is one tip: "check the controller logs." Issue #365 requests a troubleshooting guide.

7. **No decision tree** — Should I use `kelos run` or YAML? Example 01 or 02? Which agent type? Which credential type? README presents options without guidance.

---

## 3. GitHub Signals

### Repo Stats
- **Age:** ~4 weeks (created 2026-02-01)
- **Stars:** 50 | Forks: 8 | Contributors: 2
- **Community health:** 37% — no CONTRIBUTING.md, no code of conduct, no issue/PR templates
- **Open issues:** 64 (most generated by Kelos self-development pipeline)
- **Discussions:** Not enabled

### Documentation Issues (14 open, tagged `kind/docs`)

**Onboarding blockers (affects every new user):**
| Issue | Title | Priority |
|-------|-------|----------|
| #305 | Quick Start ordering creates confusion | `priority/important-longterm`, `good first issue` |
| #201 | OAuth vs API key credential confusion | — |
| #195 | Config file should link to token sources | — |

**Reference gaps (affects users past hello-world):**
| Issue | Title |
|-------|-------|
| #205 | `promptTemplate` variables incomplete (7 documented, 10 exist) |
| #221 | Valid `--model` values inconsistent |
| #202 | No examples for Codex or Gemini (`good first issue`) |
| #204 | TaskSpawner GitHub Issues repo determination unclear |

**Operational gaps:**
| Issue | Title |
|-------|-------|
| #365 | Missing troubleshooting guide |
| #242 | Missing task cleanup/lifecycle guidance |

**Use case documentation:**
| Issue | Title |
|-------|-------|
| #384 | AI-powered dependency maintenance |
| #291 | PR review, fleet migration, security audit patterns |
| #328 | Orchestrator pattern |

### CLI UX Issues (9 open)
- `kelos run` provides no next-step guidance (#183)
- `kelos run --help` doesn't mention credentials required (#162)
- `kelos run --watch` doesn't show results on completion (#402)
- Error message when controller not installed is unhelpful (#196)

### Key Insight
The `kelos-fake-user.yaml` self-development agent literally simulates a new user experience, rotating through "Documentation & Onboarding," "Developer Experience," and "Examples & Use Cases." This reveals maintainers are already aware of onboarding friction.

---

## 4. Comparable Projects

### Documentation Framework Survey

| Project | Stars | Contributors | Doc Framework | Doc Site |
|---------|-------|-------------|---------------|----------|
| Argo Workflows | 16,500 | 926 | MkDocs + ReadTheDocs | argo-workflows.readthedocs.io |
| Tekton | 8,900 | 367 | Hugo | tekton.dev/docs |
| Dagger | 15,500 | 296 | Docusaurus | docs.dagger.io |
| Concourse CI | 7,800 | ~372 | Booklit (custom) | concourse-ci.org |
| Flux CD | 7,900 | 179 | Hugo + Docsy | fluxcd.io |
| Kueue | 2,300 | — | Hugo | kueue.sigs.k8s.io |
| Crossplane | 11,500 | — | Hugo | docs.crossplane.io |

**Framework popularity in K8s ecosystem:** Hugo is the plurality choice (Tekton, Flux, Kueue, Crossplane). MkDocs used by Argo. Docusaurus used by Dagger.

### When Did They Create Doc Sites?

| Project | Stars at doc site launch | Context |
|---------|------------------------|---------|
| Tekton | ~0 (same month as launch, March 2019) | Foundation-donated project |
| Dagger | ~0 (at public launch, March 2022) | VC-funded, docs as product |
| Kueue | ~500-1,000 (early life) | kubernetes-sigs umbrella |
| Argo Workflows | ~1,000-3,000 (estimated) | Grew organically |
| Concourse | Very early | Commercial backing (Pivotal/VMware) |

**Key finding:** No project in this cohort waited until "big enough" for a doc site. Projects with institutional backing launched with doc sites from day one. Even Kueue (2,300 stars, smallest surveyed) has a full Hugo site.

### Documentation Structure Patterns

Common sections across all surveyed projects:
1. **Getting Started / Quick Start** — Every project has this
2. **Concepts** — Explaining the mental model (especially CRD-based projects)
3. **Installation / Operator Guide** — Separate from user-facing docs
4. **CLI Reference** — Auto-generated or hand-written
5. **Examples / How-to Guides** — Practical, task-oriented
6. **API Reference** — CRD field-level docs
7. **Contributing** — Developer guide

### What Dagger Does Well (most relevant model)
- Docusaurus-based (same as kelos-docsite already)
- Covers 8 SDK languages with unified structure
- Quickstarts are language-specific
- Cookbook section for practical recipes
- FAQ is prominent and addresses real user questions

---

## 5. Project Maturity Assessment

### Current State
- **Version:** v0.14.0
- **Stars:** ~50
- **Contributors:** 2 (gjkim42, aslakknutsen)
- **Age:** ~4 weeks
- **Community health:** 37%
- **Unique aspect:** Project uses itself (self-development pipeline) — a strong credibility signal

### What This Means for Documentation

**Arguments FOR a doc site now:**
1. Every comparable project created docs early (or at launch)
2. 14 open `kind/docs` issues already — documentation needs are real and growing
3. README is already long and multi-audience (users, operators, contributors)
4. The project's complexity (4 CRDs, multiple agent types, multiple credential types) exceeds what a README can handle well
5. A doc site is a credibility signal for early adopters evaluating the project

**Factors to consider:**
1. Only 2 contributors — maintenance burden matters
2. ~50 stars — external user base is tiny; most "users" are the maintainers themselves
3. APIs are still evolving (v0.14.0) — docs will need frequent updates
4. The self-development pipeline could potentially help maintain docs

### Maturity-Appropriate Documentation Strategy

At this stage, the comparable projects suggest:
- A doc site with **focused, minimal scope** — not trying to match Argo's 7-section structure
- Prioritize: Quick Start (fixing the ordering/credential issues), Concepts (the 4 CRDs), CLI Reference, Examples
- Avoid: versioned docs, multi-language support, extensive operator guides
- The existing Docusaurus setup in kelos-docsite aligns with Dagger's approach and is a reasonable framework choice

---

## 6. Existing Documentation Inventory

| Path | Contents | Audience |
|------|----------|----------|
| `README.md` | Main entry: what/why/quickstart/examples/reference/security/cost/FAQ | All |
| `docs/reference.md` | Complete field-level API reference for 4 CRDs + CLI flags | Users/operators |
| `docs/agent-image-interface.md` | Technical spec for custom agent images | Advanced users |
| `examples/README.md` | Index + how-to for all 7 examples | Users |
| `examples/NN-*/README.md` | Per-example walkthroughs | Users |
| `self-development/README.md` | Architecture + deployment guide for self-dev pipeline | Operators/contributors |

### What's Missing
- No troubleshooting guide
- No conceptual guide (when to use TaskSpawner vs `kelos run`, when to use which credential type)
- No comparison/migration guide (local Claude Code vs Kelos)
- No cost estimation guidance beyond general notes
- No CONTRIBUTING.md
- No Codex or Gemini examples despite advertised support
