# Multi-Repo Migration

Apply a repetitive change across multiple repositories using parallel agent
Tasks. This pattern is useful for fleet-wide refactoring, dependency bumps,
or API migration across microservices.

## How It Works

Each repository gets its own Workspace resource. A single TaskSpawner on
a cron schedule creates one Task per Workspace. Axon handles the
parallelism — all agents run concurrently in isolated Pods.

This example shows three repos, but the pattern scales to any number.
Add more Workspace resources and corresponding TaskSpawner entries.

Since TaskSpawner currently supports one Workspace per spawner, this
example uses one TaskSpawner per repository. A shared AgentConfig ensures
consistent behavior across all agents.

## Resources

| File | Description |
|------|-------------|
| `workspaces.yaml` | Workspace resources for each target repo |
| `credentials-secret.yaml` | Agent credentials (shared) |
| `github-token-secret.yaml` | GitHub token for cloning and pushing |
| `agentconfig.yaml` | Shared migration instructions |
| `taskspawners.yaml` | One TaskSpawner per repo, all on the same cron |

## Setup

1. Replace all `# TODO:` placeholders.
2. Add or remove Workspace/TaskSpawner pairs to match your repo list.
3. Apply:

```bash
kubectl apply -f examples/06-multi-repo-migration/
```

4. Monitor progress:

```bash
axon get tasks -w
```

## Tips

- Set `maxConcurrency: 1` on each TaskSpawner to ensure only one migration
  attempt per repo at a time.
- Use `ttlSecondsAfterFinished: 0` if you want immediate cleanup after
  each run.
- The cron schedule `"0 9 * * 1"` runs once a week (Monday 9 AM UTC).
  Adjust to run once (`"0 9 12 2 *"`) for a one-time migration.
