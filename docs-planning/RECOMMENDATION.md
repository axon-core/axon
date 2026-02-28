# Recommendation: Should Kelos Build a Documentation Site?

## Verdict: Yes, but scoped tightly — and fix the foundations first

The three analyses converge on a nuanced answer: Kelos has real documentation problems that need solving, but a full doc site isn't the first step. The right sequence is **fix the existing docs first, then build the site on a solid foundation.**

---

## Where the Analyses Agreed

All three teammates independently confirmed the same things:

1. **The existing docs are surprisingly good for a project this age.** The README, Quick Start, and 7 examples form a usable learning path. This is a strength to build on, not throw away.

2. **There are real, specific gaps.** Jira integration is completely undocumented. 8+ API fields are missing from reference.md. Credential docs are scattered across 5 locations with inconsistent naming. No troubleshooting guide exists. These aren't hypothetical — the auditor cited exact files and line numbers.

3. **Quick wins exist regardless of the doc site decision.** All three identified improvements (troubleshooting guide, reference completeness, better CLI help) that would help users today.

4. **The project's complexity already strains a README-only approach.** 4 CRDs, 4 agent types, 3 event sources, 2 auth methods, 2 credential types, dependency chaining, MCP servers, plugins — the auditor counted 83+ CRD fields and the advocate counted 15+ conceptual topics. This is too much for a single README.

## Where They Disagreed

| Question | Auditor | Researcher | Advocate |
|----------|---------|------------|----------|
| Is the README adequate? | Getting long, gaps are serious | Multi-audience problem, can't scale | "Near-complete user guide" |
| Is a doc site premature? | Implied yes (via gap severity) | No — comparable projects all launched early | Yes — 48 stars, 2 contributors |
| What's the biggest risk? | Stale/incomplete reference | Users getting stuck (14 open doc issues) | Maintenance burden on 2 people |
| What comparable projects show | — | Every surveyed project launched docs early | — |
| ROI of a doc site now | — | High (credibility signal, reduces issue volume) | Low (personal responses scale fine at 48 stars) |

**The core tension:** The researcher found that every comparable K8s project (Argo, Tekton, Dagger, Kueue, Flux) launched doc sites early — some at day one. The advocate counters that those projects had institutional backing (CNCF, VC funding, foundation donation) that Kelos doesn't. Both are right.

---

## The Recommendation

### Phase 1: Fix the foundations (do this now, ~3-5 days of work)

Before building a site, fix the content that will go in it:

1. **Complete the reference docs.** Add the 8+ missing fields (workspace `remotes`, TaskSpawner `triggerComment`/`excludeComments`/`assignee`/`author`/`priorityLabels`/`repo`, Jira source). The auditor identified every gap with line numbers — this is straightforward work.

2. **Write a troubleshooting guide** (`docs/troubleshooting.md`). All three analyses flagged this as missing. Cover: credential errors, workspace failures, rate limits, task immutability errors, pod failures, how to read controller logs.

3. **Consolidate credential documentation.** Currently scattered across 5 locations. Write one authoritative section that distinguishes agent credentials (Anthropic/OpenAI) from workspace credentials (GitHub PAT/App), and reference it from everywhere else.

4. **Add a concepts page** (`docs/concepts.md`). Answer the questions the auditor identified: When do I need a Workspace? What's the difference between AgentConfig and prompt instructions? When TaskSpawner vs manual Tasks? Include a simple decision tree.

5. **Fix the Quick Start ordering** (issue #305) and credential confusion (issue #201). The researcher identified these as blockers for every new user.

### Phase 2: Build the doc site (after Phase 1, ~1-2 weeks)

Once the content is solid, build the site with a **minimal initial scope**:

**Framework:** Docusaurus (already set up in this repo). Aligns with Dagger's approach. Good DX, React-based, versioning support when needed later.

**Initial pages — and nothing more:**

| Page | Source | Priority |
|------|--------|----------|
| Home / Overview | New (project pitch + architecture diagram) | P0 |
| Quick Start | Migrated from README, with fixes from Phase 1 | P0 |
| Concepts (4 CRDs explained) | New from Phase 1 concepts.md | P0 |
| Examples | Migrated from examples/ READMEs | P0 |
| API Reference | Migrated from docs/reference.md (completed) | P1 |
| Troubleshooting | From Phase 1 troubleshooting.md | P1 |
| Agent Image Interface | Migrated from docs/agent-image-interface.md | P2 |
| CLI Reference | Extracted from reference.md | P2 |

**Explicitly NOT in scope for v1:**
- Versioned docs (premature for v1alpha1)
- Blog
- Multi-language support
- Operator/admin guide (fold into concepts for now)
- Auto-generated API reference from Go types (nice-to-have, not now)

### Phase 3: Sustain (ongoing)

The advocate's maintenance concern is real. Mitigate it:

- **Keep docs in the same repo as code** (or require doc updates in the same PR as API changes). The advocate's key evidence — commit `628dfd5` fixing drifted example field names — shows that even co-located docs drift. A separate repo would be worse.
- **Use CI to check for drift.** A simple script that compares documented fields against Go struct tags would catch the reference gap problem automatically.
- **Leverage the self-development pipeline.** The researcher noted that `kelos-fake-user.yaml` already simulates new user experiences. A doc-review agent could flag stale content.
- **Resist scope creep.** The advocate's "empty restaurant" warning is valid — 8 solid pages beat 20 stub pages.

---

## Signals That Would Change This Recommendation

If Phase 1 fixes resolve most user friction (tracked via GitHub issue closure rate), you could delay Phase 2 longer. Conversely, accelerate Phase 2 if:

- Stars cross 200+ (the advocate's threshold)
- More than 2-3 external contributors appear
- GitHub issues tagged `kind/docs` keep growing despite Phase 1 fixes
- Users explicitly ask for a doc site

---

## Quick Wins — Do These Regardless

These cost hours, not days, and help immediately:

1. **Document Jira integration.** It's fully coded and completely invisible. Even a section in the README unblocks Jira users.
2. **Add missing fields to reference.md.** The auditor listed them all with line numbers.
3. **Deduplicate the template variables table.** Currently in 3 places (reference.md, self-development/README.md, partially in README). Single source of truth, link from elsewhere.
4. **Add `CONTRIBUTING.md`.** Community health is at 37%. This is table stakes.
5. **Fix the credential naming inconsistency.** `token` vs `secretRef` vs `oauthToken` — add a mapping table or unify the terminology.

---

## Summary

The auditor showed the gaps are real and specific. The researcher showed that comparable projects don't wait. The advocate showed that maintenance burden on 2 contributors is a genuine constraint. The synthesis: **do both, in sequence.** Fix the content first (Phase 1), then build the site on solid foundations (Phase 2). Don't let perfect be the enemy of good — but also don't build a house on a cracked foundation.

The fact that `kelos-docsite` already exists as a repo suggests you've already been thinking about this. The recommendation is: finish the content work before investing heavily in the presentation layer.
