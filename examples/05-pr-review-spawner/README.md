# PR Review Spawner

Automatically review open pull requests using an AI agent. The spawner
discovers PRs with a specific label and creates a Task for each one. The
agent reads the diff, checks for common issues, and posts a review comment.

## How It Works

1. A developer opens a PR and adds the `needs-ai-review` label.
2. TaskSpawner discovers the labeled PR and creates a Task.
3. The agent clones the repo, reads the PR diff and description, and posts
   a review via the `gh` CLI.
4. After the review, the agent removes the trigger label and adds
   `ai-reviewed` so it is not picked up again.

## Resources

| File | Description |
|------|-------------|
| `workspace.yaml` | Git repo the agent clones |
| `credentials-secret.yaml` | Agent credentials (OAuth token) |
| `github-token-secret.yaml` | GitHub token for cloning and `gh` CLI |
| `agentconfig.yaml` | Review instructions and skill |
| `taskspawner.yaml` | TaskSpawner watching for labeled PRs |

## Setup

1. Replace all `# TODO:` placeholders in the YAML files.
2. Create the GitHub labels `needs-ai-review` and `ai-reviewed` in your repo.
3. Apply all resources:

```bash
kubectl apply -f examples/05-pr-review-spawner/
```

4. Label a PR with `needs-ai-review` and watch the agent pick it up:

```bash
axon get tasks -w
```
