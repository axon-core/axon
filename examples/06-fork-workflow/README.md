# 06 — Fork Workflow

A TaskSpawner that discovers issues from an **upstream** repository but runs
agents against a **fork**. This is the standard open-source contribution
pattern where you don't have push access to the upstream repo.

## Use Case

You maintain a fork of an upstream project. When new issues are filed upstream,
an AI agent clones your fork, creates a fix branch, and opens a PR from your
fork back to the upstream repository.

## Resources

| File | Kind | Purpose |
|------|------|---------|
| `workspace.yaml` | Workspace | Points to your fork with an `upstream` remote |
| `taskspawner.yaml` | TaskSpawner | Polls upstream for issues via `githubIssues.repo` |

You also need Secrets for GitHub and agent credentials (see example 03).

## How It Works

```
TaskSpawner polls upstream repo for issues (via githubIssues.repo)
    │
    ├── new issue found → creates Task
    │       │
    │       ├── clones fork (origin)
    │       ├── adds upstream remote
    │       ├── creates fix branch
    │       ├── agent works on fix
    │       └── opens PR from fork → upstream
    └── ...
```

## Key Concepts

- **`workspace.spec.repo`** — points to your fork (this is what gets cloned
  as `origin`).
- **`workspace.spec.remotes`** — adds the upstream repo as a named remote so
  the agent can reference it.
- **`taskspawner.spec.when.githubIssues.repo`** — tells the spawner to poll
  the upstream repo for issues instead of the fork.
- The same GitHub token (from `workspace.spec.secretRef`) is used for both
  polling upstream issues and pushing to the fork.

## Steps

1. **Create secrets** for your GitHub token and agent credentials.

2. **Edit `workspace.yaml`** — replace the fork and upstream URLs.

3. **Edit `taskspawner.yaml`** — set `githubIssues.repo` to the upstream repo.

4. **Apply the resources:**

```bash
kubectl apply -f examples/06-fork-workflow/
```

5. **Watch for spawned Tasks:**

```bash
kubectl get tasks -w
```

6. **Cleanup:**

```bash
kubectl delete -f examples/06-fork-workflow/
```
