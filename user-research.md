# User Perspective Research — Kelos Documentation Site Analysis

## 1. GitHub Issues & Discussions Signal

### Volume: 38 total `kind/docs` issues (14 open, 24 closed)

For a project that is **less than 1 month old** with only **50 stars** and **2 contributors**, having 38 documentation-specific issues is a massive signal. That's roughly 1 docs issue for every 1.3 stars.

### Open Docs Issues (14) — Categorized

**Onboarding/Quick Start friction (5 issues):**
- #305: Quick Start ordering creates confusion about when to install
- #201: OAuth vs API key credential type is confusing
- #195: Config file should link to where users can obtain tokens
- #365: Missing troubleshooting guide for common failures
- #221: Missing guidance on valid model values for --model flag

**Missing documentation for existing features (5 issues):**
- #242: Missing guidance on task cleanup and lifecycle management
- #205: TaskSpawner promptTemplate variables incomplete
- #204: TaskSpawner GitHub Issues example doesn't explain how repo is determined
- #202: Missing usage examples for Codex and Gemini agent types
- #207: Self-development TaskSpawners use aggressive 1m pollInterval without explanation

**Missing guides/use cases (3 issues):**
- #384: New use case: AI-powered dependency maintenance
- #291: New use cases: PR review, fleet migration, security audit, and beyond
- #328: Orchestrator pattern

**Meta (1 issue):**
- #219: Update contributing guide

### Closed Docs Issues (24) — Patterns

Many closed issues reveal problems users already hit:
- #454: Example 05 had incorrect field names (broken YAML)
- #436: GitHub App auth was fully implemented but completely undocumented
- #243: Quick Start didn't explain where to obtain credentials
- #229: Quick Start lacked context about what happens during first run
- #220: No explanation of where task output goes without a workspace
- #251: Demo used misleading custom task name without explaining --name
- #222, #231: Missing standalone YAML examples (filed twice — strong signal)
- #245: Inconsistent secret naming across examples
- #159, #163, #161: Multiple Quick Start gaps fixed

### Discussions

Only 1 discussion exists (the default welcome post, #288). No community Q&A activity yet.

### Key Takeaway

The docs issues represent **real user friction**, not feature requests. Multiple issues describe the exact moment where a new user gets stuck: credential setup, understanding task output, debugging failures. These are the kinds of problems a doc site solves better than a README.

---

## 2. Likely Users — Who Would Use Kelos?

### Primary Persona: Platform Engineers / DevOps Engineers
- **Background:** Kubernetes-fluent, familiar with operators and CRDs, comfortable with `kubectl` and YAML
- **Motivation:** Wants to automate repetitive coding tasks at scale — fix bugs, update docs, refactor across repos
- **Discovery path:** Likely finds Kelos through Hacker News, Twitter/X, or searching for "Kubernetes AI agent orchestration"
- **Pain points:** Already juggling many tools; needs Kelos to be immediately understandable or they'll move on

### Secondary Persona: AI/ML Engineers Exploring Agentic Workflows
- **Background:** Strong on AI/ML, moderate Kubernetes knowledge
- **Motivation:** Wants to run Claude/Codex/Gemini agents at scale without building infra
- **Discovery path:** Finds via "Claude Code Kubernetes" or "autonomous AI agents" searches
- **Pain points:** May need more hand-holding on Kubernetes concepts; the current docs assume K8s fluency

### Tertiary Persona: Engineering Managers / Tech Leads
- **Background:** Evaluating tools for team adoption
- **Motivation:** Looking for ROI story — can this automate issue triage, PR review, CI fix?
- **Discovery path:** Team member shares it, or sees it mentioned in AI/DevOps newsletter
- **Pain points:** Needs to quickly understand value prop, security model, and cost implications without reading all the YAML

### What All Personas Need
1. **Quick understanding** of what Kelos does (README handles this well)
2. **Step-by-step getting started** that doesn't assume prior context (current Quick Start has gaps)
3. **Reference docs** they can search/navigate when building real workflows (currently a single reference.md file)
4. **Patterns and recipes** beyond hello-world (partially covered by examples/)
5. **Troubleshooting** when things go wrong (currently missing)

---

## 3. User Journey Map

### Step 1: Discovery (README.md)
**Experience:** Good. The README has a clear value prop, demo CLI output, and a video. The "Why Kelos?" section is compelling.
**Friction:** Moderate. Users see the demo but can't try it immediately — they need a K8s cluster, credentials, and setup first.

### Step 2: Decide to Try It (Prerequisites)
**Experience:** OK. Prerequisites are listed (K8s 1.28+), with a `kind` expandable section.
**Friction:** High for non-K8s users. No guidance on what "Kubernetes cluster" means practically or how much it costs.

### Step 3: Install CLI
**Experience:** Simple — one curl command or `go install`.
**Friction:** Low.

### Step 4: Install Kelos to Cluster
**Experience:** One command (`kelos install`).
**Friction:** Medium. Issue #305 notes that the ordering (install before init) is confusing. The `kelos init` output says "Install Kelos (if not already installed)" suggesting init should come first.

### Step 5: Configure Credentials (~/.kelos/config.yaml)
**Experience:** This is the biggest pain point.
**Friction:** Very High.
- Issue #201: OAuth vs API key confusion — README shows OAuth first, but CLI defaults to api-key
- Issue #195: Config file doesn't link to where to get tokens
- Users need BOTH an AI provider credential AND optionally a GitHub token — this two-credential requirement isn't clearly communicated

### Step 6: Run First Task
**Experience:** The command itself is simple (`kelos run -p "..."`)
**Friction:** Medium.
- Issue #220 (closed): Users didn't understand that without a workspace, output is ephemeral
- Issue #163 (closed): Users didn't know how to get the task name from output

### Step 7: Debug When Things Go Wrong
**Experience:** Very poor.
**Friction:** Very High.
- Only guidance: "check controller logs" (one line in README)
- No troubleshooting section, no common failure patterns
- Users are left guessing whether it's a credential issue, cluster issue, or agent issue

### Step 8: Build Real Workflows (Beyond Hello World)
**Experience:** Moderate. 7 examples exist covering common patterns.
**Friction:** Medium-High.
- Examples are YAML-heavy with `# TODO:` placeholders
- No narrative explanation of WHY you'd use each pattern
- Reference doc is a single flat file — no search, no deep linking, no navigation

### Journey Score Card

| Step | Friction | Current Doc Coverage |
|------|----------|---------------------|
| Discovery | Low | README (good) |
| Prerequisites | Medium | README (adequate) |
| Install CLI | Low | README (good) |
| Install Kelos | Medium | README (ordering issue) |
| Configure Credentials | **Very High** | README (improved, still confusing) |
| Run First Task | Medium | README (adequate) |
| Debug Failures | **Very High** | Almost none |
| Build Real Workflows | Medium-High | Examples + reference.md |

---

## 4. Comparable Projects Analysis

### Argo Workflows
- **Stars:** 16,483
- **Created:** Aug 2017 (by Applatix, acquired by Intuit 2018)
- **Doc site:** Yes — [argoproj.github.io/workflows](https://argoproj.github.io/workflows/) and [argo-workflows.readthedocs.io](https://argo-workflows.readthedocs.io/en/latest/)
- **When:** Docs evolved with the project. Initially in-repo markdown, then moved to mkdocs + ReadTheDocs as the project grew. The GitHub Pages umbrella site serves all Argo projects.
- **CNCF:** Incubating April 2020, Graduated March 2022
- **Doc structure:** Getting Started, Walk Through, Examples, Fields Reference, CLI Reference, Architecture
- **Key insight:** Argo started with in-repo markdown and examples. The formal doc site came as the project reached hundreds of stars and CNCF adoption. They use mkdocs (docs-as-code) which keeps docs in the same repo.

### Tekton
- **Stars:** 8,901
- **Created:** Aug 2018 (spun off from Knative Build)
- **Doc site:** Yes — [tekton.dev](https://tekton.dev/docs/)
- **When:** tekton.dev launched early in the project's lifecycle, around the time of CDF donation (March 2019). Beta release of the Pipeline API in March 2020.
- **CDF:** Donated to Continuous Delivery Foundation March 2019, Graduated October 2022
- **Doc structure:** Getting Started (Tasks, Pipelines), Concepts, How-to Guides, Reference, Blog
- **Key insight:** Tekton had organizational backing (Google, CDF) from the start, which funded proper docs infrastructure early. Their docs follow the [Diátaxis framework](https://diataxis.fr/) (tutorials, how-to guides, reference, explanation).

### Dagger
- **Stars:** 15,477
- **Created:** Nov 2019 (by Solomon Hykes, Docker founder)
- **Doc site:** Yes — [docs.dagger.io](https://docs.dagger.io/)
- **When:** Launched simultaneously with the public beta in March 2022. Docs were ready on day one.
- **Funding:** $20M Series A at launch
- **Doc structure:** Quickstarts, Key Concepts, Guides, Examples, API Reference, Cookbook, FAQ
- **Key insight:** Dagger launched with docs-as-a-first-class-product. Being VC-backed, they invested in docs before going public. Their docs site went through a major redesign when they pivoted from CUE to SDK-based approach.

### Summary Comparison

| Project | Stars at doc site launch | Doc site timing | Key driver |
|---------|------------------------|-----------------|------------|
| Argo Workflows | ~500-1000 (est.) | Evolved gradually from in-repo markdown | Community growth, CNCF process |
| Tekton | ~100-500 (est.) | Very early (org-backed) | Google/CDF backing, enterprise adoption |
| Dagger | ~0 (launched together) | Day 1 of public launch | VC funding, product-company model |
| **Kelos** | **50** | **None yet** | **Early-stage, 2 contributors** |

### What Kelos Can Learn

1. **All comparable projects have doc sites** — this is table stakes for Kubernetes-native tooling targeting platform engineers.
2. **docs-as-code (mkdocs/docusaurus) is the standard** — Argo uses mkdocs, Tekton uses a custom site, Dagger uses Docusaurus. All keep docs in the same repo as code.
3. **Start simple, grow with the project** — Argo's approach of starting with in-repo markdown and evolving to a doc site is the most applicable model for Kelos's stage.
4. **The Diátaxis framework works** — Tekton's structure (tutorials, how-to, reference, explanation) is the gold standard for K8s tooling docs.
5. **A doc site solves the "single README" problem** — At 584 lines, the Kelos README is already dense. A doc site enables search, navigation, deep linking, and progressive disclosure.

---

## 5. Repo Traffic & Growth Context

- **490 unique visitors** in the last 2 weeks (Feb 14-27)
- **3,943 total views** in the same period
- **Peak day:** Feb 27 with 574 views (41 unique)
- **Average:** ~35 unique visitors/day

This traffic level means real people are discovering and evaluating Kelos every day. Each visitor hitting the README wall on credentials or debugging is a potential user lost.

---

## 6. Synthesis

### The Case FOR a Doc Site (from user perspective)

1. **38 docs issues in <1 month** = users are struggling with documentation
2. **The user journey has 2 "very high friction" points** (credentials, debugging) that a README can't solve well
3. **All comparable K8s-native workflow tools have doc sites** — it's an expectation for this audience
4. **490 unique visitors in 2 weeks** — real users are evaluating Kelos right now
5. **The README is at 584 lines** and growing — it's past the point where a single file serves well
6. **Search and navigation** are impossible in the current flat-file structure

### The Case AGAINST (or "not yet")

1. **50 stars, 2 contributors** — the content isn't stabilized yet; a doc site adds maintenance overhead
2. **Argo's path** (gradual evolution from in-repo markdown) is also valid and lower-effort
3. **The docs content needs to be written first** — a doc site is just infrastructure; without content improvements, it's lipstick on a pig
4. **Opportunity cost** — time spent on doc site infrastructure is time not spent on features

### Recommendation from User Perspective

A documentation site would **meaningfully improve the user experience** at this stage. The current pain points (credential confusion, missing troubleshooting, undiscoverable reference docs) are best solved by structured documentation with navigation, search, and progressive disclosure.

The minimum viable approach: use **mkdocs-material** (like Argo) or **Docusaurus** (like Dagger), keep docs in the same repo, and start with the content that addresses the top user friction points:
1. Improved Getting Started tutorial
2. Credential/authentication guide
3. Troubleshooting guide
4. Searchable API reference
