# 08 — Task with Kelos Skill

Run a Task whose agent knows how to author and debug Kelos resources, powered
by the first-party Kelos skill.

## What This Demonstrates

- Injecting the Kelos skill via an AgentConfig plugin
- Giving the agent knowledge of Kelos CRDs, CLI, and troubleshooting

## Resources

| File | Resource | Purpose |
|------|----------|---------|
| `agentconfig.yaml` | AgentConfig | Defines the `kelos` plugin with the Kelos skill |
| `task.yaml` | Task | Runs an agent that uses the skill to set up a Kelos workflow |

## Prerequisites

1. A running Kubernetes cluster with Kelos installed (`kelos install`)
2. A Workspace resource named `my-workspace` pointing to your repo
3. A Secret named `claude-oauth-token` with OAuth credentials

## Usage

```bash
# Create the AgentConfig
kubectl apply -f examples/08-task-with-kelos-skill/agentconfig.yaml

# Run the Task
kubectl apply -f examples/08-task-with-kelos-skill/task.yaml

# Watch progress
kubectl get tasks -w
```

## Customizing

The skill content in `agentconfig.yaml` is a condensed version. For the full
skill with complete reference YAML patterns, see `skills/kelos/SKILL.md` and
`skills/kelos/references/` in the repository root.

To use the full skill content, replace the inline `content` with the contents
of `skills/kelos/SKILL.md`, or use the CLI:

```bash
kelos create agentconfig kelos-skill-agent \
  --skill kelos=@skills/kelos/SKILL.md
```
