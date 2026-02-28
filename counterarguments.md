# The Case Against Building a Documentation Site Right Now

## Executive Summary

Kelos is a **27-day-old project** with **one dominant contributor** (99.6% of commits), a **v1alpha1 API with no stability guarantees**, and **50+ open issues** — most of which are feature requests still awaiting triage. Building a documentation site now would be premature, expensive to maintain, and solve the wrong problem. The real documentation gaps are a missing troubleshooting guide, absent CONTRIBUTING.md, and weak CLI help text — none of which require a dedicated site.

---

## 1. Premature Optimization

### The numbers tell a clear story

| Metric | Value | What it means |
|--------|-------|---------------|
| Project age | 27 days (created Feb 1, 2026) | Barely past initial development |
| GitHub stars | 50 | Minimal external awareness |
| Forks | 8 | Very few external users |
| Contributors | 3 (548 of 550 commits from one person) | Single-author project in practice |
| Version | v0.14.0 (pre-1.0) | Not yet considered stable |
| Open issues | 50+ (40+ marked `kelos/needs-input`) | Core design decisions still pending |
| Release cadence | ~1 release every 2-3 days | API surface expanding rapidly |

### What this means

The project is in the **"feature explosion" phase** — capabilities are being added faster than they can be documented. In the last 27 days, Kelos has gone from v0.4.0 to v0.14.0 (10 minor versions), adding CronJob support, Jira integration, MCP servers, GitHub App authentication, TaskSpawner filtering, and task pipelines.

A documentation site assumes a **stable surface to document**. Kelos doesn't have that yet. The 40+ issues marked `kelos/needs-input` represent unresolved design decisions that could reshape the API. Building a doc site now means building on sand.

### Comparable project timing

Most Kubernetes-native projects didn't build dedicated doc sites until they had:
- Hundreds of stars and dozens of contributors
- A stable (v1+) API or at least a beta API
- External production users generating real documentation needs
- A contributor community that needed onboarding documentation

Kelos has none of these yet.

---

## 2. Maintenance Cost Will Be Brutal

### API surface area

Kelos has **4 CRDs** with substantial complexity:

| CRD | JSON fields | Validation rules | Nested struct types |
|-----|------------|------------------|---------------------|
| Task | 33 | 8 | 5 |
| TaskSpawner | 51 | 16 | 8+ |
| Workspace | 16 | 7 | 3 |
| AgentConfig | 24 | 5 | 5 |
| **Total** | **124 fields** | **36 validation rules** | **20+ nested types** |

Every one of those 124 fields needs accurate documentation. Every validation rule needs to be explained. Every nested type needs its own section.

### API change velocity

The `api/v1alpha1/` directory has seen **57 commits** in the project's 27-day lifetime — roughly **2 API changes per day**. Recent changes include:

- Breaking field relocations (workspaceRef moved from GitHubIssues to TaskTemplate level)
- New nested structures (secretRef field nesting for headers/env)
- New CRD capabilities (MCP servers, Jira integration, GitHub App auth)
- Validation rule changes (CEL rules, enum constraints, immutability enforcement)
- New required fields added to existing types

### The maintenance math

At the current rate of ~2 API changes per day:
- A doc site would need **daily review** to stay accurate
- Each breaking change requires updating reference docs, examples, and potentially tutorials
- With a single contributor doing 99.6% of development, **documentation maintenance directly competes with feature development**

The v1alpha1 designation means there are **no stability guarantees**. Fields can be removed, renamed, or restructured in any release. A doc site documenting today's API could be materially wrong within a week.

### What stale docs cost

Stale documentation is **worse than no documentation**. A user who follows an outdated doc site tutorial and gets cryptic errors will:
1. Lose trust in the project
2. File issues that waste developer time
3. Conclude the project is abandoned or poorly maintained

With the README, staleness is obvious — it's right there in the repo, updated alongside the code. A separate doc site creates a second source of truth that silently drifts.

---

## 3. Better Ways to Spend the Same Effort

The existing documentation is actually quite good for a 27-day-old project:

- **README.md**: 583 lines, well-structured with Quick Start, architecture diagrams, 5 example patterns, security considerations, cost/limits section, and FAQ
- **docs/reference.md**: 309 lines of field-by-field CRD reference
- **docs/agent-image-interface.md**: 132 lines of agent interface specification
- **examples/**: 7 complete working examples with READMEs covering all major use cases
- **self-development/README.md**: Excellent workflow documentation with diagrams

The effort to build and maintain a doc site (conservatively 20-40 hours to set up, plus ongoing maintenance) could instead address the **actual gaps** in the current documentation:

### Gap 1: No troubleshooting guide (HIGH impact, ~2 hours)

There is no troubleshooting section anywhere in the project. Users who hit errors have no guidance. Common scenarios that need coverage:
- TaskSpawner creates no tasks (wrong labels, missing secret)
- Task fails immediately (invalid credentials, kubeconfig issues)
- Task hangs in Waiting phase (dependency name typo)
- `kelos run` reports "secret not found"

### Gap 2: No CONTRIBUTING.md (HIGH impact, ~1 hour)

The project has no contributor guide at all. For a project that wants to grow beyond one contributor, this is a critical gap. It should cover:
- Local dev setup (`make verify`, `make test`)
- Code organization (cmd/, internal/, api/)
- PR expectations and review process

### Gap 3: CLI help text lacks examples (MEDIUM impact, ~2 hours)

The CLI commands have `Short` descriptions but no `Long` descriptions or `Example` fields. Users running `kelos run --help` see a description but no practical usage example. Adding examples to the cobra commands is trivial and immediately helpful.

### Gap 4: Error messages lack actionable guidance (MEDIUM impact, ~3 hours)

Error messages use proper wrapping (`fmt.Errorf` with `%w`) but don't tell users what to do next. For example:
- "creating discovery client: %w" — should suggest checking kubeconfig
- Credential failures — should suggest which auth method to try
- Workspace not found — should suggest `kelos get workspaces`

### Gap 5: Missing examples for edge cases (LOW impact, ~3 hours)

The 7 existing examples cover the happy path well. Missing scenarios:
- Resource constraints (CPU/memory limits)
- Private repository authentication (SSH keys, GitHub App)
- Error recovery and debugging workflows

**Total cost of all five gaps: ~11 hours.** That's less than half the effort of setting up a doc site, and addresses the problems users actually have.

---

## 4. Risks of Building Too Early

### Risk 1: Version drift between docs and code

A doc site creates a **second source of truth**. The README lives in the repo and is updated in the same PR that changes the code. A doc site requires a separate update process. With 2 API changes per day and one contributor, this synchronization will break.

### Risk 2: Splitting attention from core development

The project has 50+ open issues, most marked `kelos/needs-input`. The sole active contributor's time is the project's most scarce resource. Every hour spent on doc site infrastructure is an hour not spent on:
- Resolving the 40+ issues awaiting design decisions
- Stabilizing the API toward v1beta1
- Adding tests (the project has integration and e2e tests, but coverage could improve)
- Improving error messages and CLI UX

### Risk 3: Premature information architecture

A doc site imposes structure — navigation, categories, hierarchy. With the API still in flux, that structure will need repeated reorganization. Today's "Getting Started" tutorial may reference fields that don't exist next month. Today's "Advanced Configuration" page may become basic default behavior.

### Risk 4: False signal of maturity

A polished doc site signals to users that the project is stable and production-ready. Kelos is v0.14.0 with a v1alpha1 API. Users who see a professional doc site may deploy to production and be surprised by breaking changes. The current README correctly sets expectations — it's clearly a project README, not a product manual.

### Risk 5: Documentation debt accumulation

Once a doc site exists, every feature must be documented there. This creates a new category of work (documentation debt) that competes with code debt. For a single-contributor project, this can become a psychological burden that slows development — "I can't ship this feature until I update the docs" or worse, "I shipped but forgot to update the docs."

---

## 5. Minimum Viable Alternative

Instead of a documentation site, here's the smallest set of improvements that would meaningfully help users:

### Tier 1: Do this week (~4 hours total)

1. **Add a Troubleshooting section to README.md** — 10-15 FAQ entries covering the most common failure modes with solutions. This is the single highest-impact documentation improvement possible.

2. **Create CONTRIBUTING.md** — Local dev setup, test commands, code organization, PR guidelines. Unblocks community contributions.

3. **Add `Example` fields to CLI cobra commands** — Users running `--help` see practical usage patterns immediately.

### Tier 2: Do this month (~4 hours total)

4. **Improve error messages** — Add actionable next steps to the 10 most common error paths. This is documentation embedded in the product itself — it can never go stale.

5. **Add a "What to do when things go wrong" section to examples/** — Each example README gets 5-10 lines about common failures and debugging.

### Tier 3: Do when the API stabilizes

6. **Consider a doc site** — Once the API reaches v1beta1 or v1, the project has 10+ contributors, and there's a genuine user community generating documentation needs, a doc site becomes justified.

### Why this works

- All improvements live **in the repo**, so they're versioned with the code
- No new infrastructure to maintain (no static site generator, hosting, CI pipeline)
- Each improvement is **independently valuable** — you can ship one without the others
- Total effort: ~8 hours for Tiers 1-2, delivering more user value than a doc site would

---

## Conclusion

The case against a doc site is not that documentation doesn't matter — it absolutely does. The case is that **a doc site is the wrong form factor for this project at this stage**. The project is 27 days old, has one active contributor, a rapidly changing v1alpha1 API, and 50+ unresolved design decisions. The actual documentation gaps (troubleshooting, contributor guide, CLI help, error messages) are best addressed with targeted, in-repo improvements that stay in sync with the code by default.

Build a doc site when:
- The API reaches v1beta1 or later
- There are 5+ regular contributors who need onboarding
- There are 200+ stars indicating genuine external user base
- Users are filing issues that can only be solved with structured documentation (tutorials, conceptual guides, migration paths)

Until then, invest in making the existing README, examples, CLI help, and error messages as good as they can be. That's where users actually look for help, and that's where improvements will have the highest return.
