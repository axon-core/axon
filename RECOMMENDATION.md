# Documentation Site Recommendation

## Verdict: Not Yet — But Sooner Than You'd Think

The answer is **not yet**, but with a concrete trigger: **build the doc site when the API reaches v1beta1 or when you're ready to actively recruit contributors/users** — whichever comes first. In the meantime, there are high-impact improvements that take hours, not days, and address the real pain users are hitting right now.

This isn't a "no" — it's a "sequence it right."

---

## How We Got Here

Three independent analyses examined this from different angles:

| Analyst | Perspective | Conclusion |
|---------|------------|------------|
| **Auditor** ([docs-audit.md](docs-audit.md)) | What do we actually have? | Docs are strong for the stage but hitting limits — 3 major features undocumented, no conceptual guides, poor discoverability |
| **Researcher** ([user-research.md](user-research.md)) | What do users experience? | 38 docs issues in <1 month, 2 "very high friction" points in the user journey, all comparable projects have doc sites |
| **Devil's Advocate** ([counterarguments.md](counterarguments.md)) | Why not build it now? | 27-day-old project, 1 real contributor, 2 API changes/day, v1alpha1 instability, targeted fixes are cheaper |

---

## Where They Agreed

All three analyses converged on these points:

1. **The existing README is genuinely good.** Well-structured, honest, progressive disclosure via `<details>` blocks. This is not a documentation-neglected project.

2. **Troubleshooting documentation is the #1 gap.** The auditor found no troubleshooting guide anywhere. The researcher identified debugging as a "very high friction" journey step. The devil's advocate lists it as the highest-impact fix (~2 hours).

3. **Several shipped features have zero documentation.** Jira integration, comment-based triggers (`TriggerComment`/`ExcludeComments`), priority labels, assignee/author filters — all exist in the API types but appear nowhere in docs. This is a content problem, not a format problem.

4. **Credential setup is confusing.** Five separate GitHub issues about OAuth vs API key, where to get tokens, and which credentials are needed when. The auditor found four different explanations of workspace auth scattered across files.

5. **The content must improve regardless of format.** A doc site with the same gaps is just a fancier version of the same problem. Content first, infrastructure second.

---

## Where They Disagreed

### The timing question

The researcher argues **now**: 38 docs issues in a month, 490 unique visitors in 2 weeks, all comparable K8s tools have doc sites — users expect it and are hitting walls without it.

The devil's advocate argues **not yet**: 27 days old, 1 contributor doing 99.6% of work, 2 API changes per day, v1alpha1 means no stability guarantees. A doc site is premature optimization.

**My read:** Both are right about different things. The user friction is real and documented (38 issues isn't noise). But the API instability is also real — documenting 124 fields that change twice daily is a maintenance trap. The resolution is to **fix the content problems now** (in-repo, where they stay in sync) and **build the site infrastructure when the API stabilizes** (so you're not maintaining two sources of truth during a period of rapid change).

### The form factor question

The researcher says a 584-line README can't scale and needs search/navigation. The devil's advocate says in-repo markdown is the right format because it versions with the code.

**My read:** These aren't mutually exclusive. Modern doc site tools (mkdocs-material, Docusaurus) render in-repo markdown into a searchable site. The content stays in the repo, versioned with the code. The "site" is just a build step. But setting up that build step has a cost (CI pipeline, hosting, theme config, nav structure) that isn't justified until the content is worth navigating.

---

## Recommendation

### Phase 1: Quick wins now (this week, ~8 hours)

These improvements address real user pain, require no infrastructure, and stay in sync with code:

| Action | Impact | Effort | Addresses |
|--------|--------|--------|-----------|
| Add troubleshooting FAQ to README | High | ~2 hrs | #365, user journey step 7 |
| Create CONTRIBUTING.md | High | ~1 hr | #219, contributor onboarding |
| Document Jira integration in reference.md | High | ~1 hr | Completely missing feature |
| Document TriggerComment, ExcludeComments, PriorityLabels, Assignee, Author in reference.md | High | ~1 hr | 5 undocumented fields |
| Add Workspace `remotes` to reference.md | Medium | ~30 min | Missing from reference, only in example 06 |
| Add `Example` blocks to CLI cobra commands | Medium | ~2 hrs | CLI discoverability |
| Deduplicate or differentiate AGENTS.md vs CLAUDE.md | Low | ~30 min | Identical files causing confusion |

### Phase 2: Content expansion (this month, ~6 hours)

Write the missing conceptual content — still as in-repo markdown, but structured so it can become doc site pages later:

| Action | Notes |
|--------|-------|
| Credential & authentication guide | Unified guide: API key vs OAuth vs GitHub App, which secrets for what, end-to-end flow |
| Task lifecycle guide | Phases, branch serialization, dependency resolution, with a state diagram |
| Architecture overview | Controller, Jobs/Deployments/CronJobs under the hood, reconciliation loop |
| Improve error messages in code | Actionable guidance in the errors themselves — this is documentation that can never go stale |

Put these in `docs/` with a `docs/README.md` index. This creates a natural structure that maps directly to doc site pages later.

### Phase 3: Doc site (when ready)

**Trigger criteria** — build the site when at least 2 of these are true:
- API reaches v1beta1 (stability guarantee means docs won't churn daily)
- 5+ regular contributors (onboarding cost justifies the investment)
- 200+ stars (enough users that discoverability matters at scale)
- You're actively promoting the project (conference talks, blog posts, outreach)

**When you build it:**
- **Framework:** mkdocs-material. It's what Argo Workflows uses, it renders in-repo markdown, has search built in, and is the standard for K8s ecosystem projects. Docusaurus is also viable but is heavier (React-based) and more common in the JS ecosystem.
- **Structure (Diataxis framework):**
  - Tutorials: Getting Started, First TaskSpawner, First Pipeline
  - How-to Guides: Authentication, Fork Workflows, Custom Agent Images, Jira Integration
  - Reference: CRD specs (auto-generated from types if possible), CLI reference, Config file reference
  - Explanation: Architecture, Task Lifecycle, Security Model
- **Scope:** Start with the content from Phases 1-2 reorganized into the structure above. Don't write new content just for the site launch.
- **Keep it in-repo:** `docs/` directory, built by CI, deployed to GitHub Pages. One repo, one PR process, one review cycle.

---

## What NOT to Do

1. **Don't build a doc site and leave the content gaps.** A searchable site with no troubleshooting guide is still useless for a stuck user.
2. **Don't auto-generate CRD reference docs from Go types right now.** The types change too fast. Wait for v1beta1 when the structure stabilizes.
3. **Don't split docs into a separate repo.** Version drift is the #1 risk the devil's advocate identified, and a separate repo maximizes that risk.
4. **Don't add a blog/changelog to the doc site.** GitHub Releases already serves this purpose. Don't create another thing to maintain.

---

## Summary

The project has strong documentation for its age. The pain points are specific and addressable: missing troubleshooting, undocumented features, credential confusion. These are content problems that should be fixed in-repo now. A doc site is the right eventual destination — it's what every comparable project has — but building one before the API stabilizes and the content gaps are filled would create maintenance burden without proportional user value.

Do the content work first. The site infrastructure is a couple days of work when the time is right. The content is the hard part, and you can start that today.
