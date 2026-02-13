# Security Audit Cron

Run a periodic security audit of your codebase using an AI agent. The agent
checks for common vulnerabilities, outdated dependencies with known CVEs,
hardcoded secrets, and insecure patterns. Results are filed as GitHub issues.

## How It Works

1. A cron schedule triggers the audit (default: weekly on Monday mornings).
2. The agent clones the repo, scans for security issues, and checks
   dependencies.
3. If issues are found, the agent creates a GitHub issue per finding (or one
   summary issue) with remediation guidance.
4. If no issues are found, the agent exits without creating any output.

## Resources

| File | Description |
|------|-------------|
| `workspace.yaml` | Git repo to audit |
| `credentials-secret.yaml` | Agent credentials |
| `github-token-secret.yaml` | GitHub token for creating issues |
| `agentconfig.yaml` | Security audit instructions |
| `taskspawner.yaml` | Weekly cron TaskSpawner |

## Setup

1. Replace all `# TODO:` placeholders.
2. Apply:

```bash
kubectl apply -f examples/07-security-audit-cron/
```

3. Check results after the next cron tick:

```bash
axon get tasks
axon logs <task-name>
```

## Customization

- Adjust the cron schedule in `taskspawner.yaml` to match your needs.
- Edit the audit checklist in `agentconfig.yaml` to add project-specific
  security requirements.
- Set `maxConcurrency: 1` to prevent overlapping audit runs.
