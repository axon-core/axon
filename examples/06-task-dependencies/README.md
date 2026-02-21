# 06 — Task Dependencies (Orchestrator Pattern)

Chain multiple Tasks together so one agent's output feeds into the next.
This pattern lets you break complex work into specialized steps — for
example, one agent writes code and another reviews it.

## Use Case

Run a multi-step workflow where a coding agent creates a PR, then a
review agent examines the PR and leaves feedback. The review agent
automatically receives the branch name and PR URL produced by the first
agent through prompt templates.

## How It Works

Axon's `dependsOn` field keeps Task B in a `Waiting` phase until Task A
succeeds. When Task A completes, Axon captures its structured outputs
(lines matching `key: value` printed between output markers) and makes
them available to Task B via Go template syntax in the prompt:

```
{{ index .Deps "<task-name>" "Results" "<key>" }}
```

The controller resolves the template before creating Task B's Job, so the
downstream agent sees a fully rendered prompt with concrete values.

## Resources

| File | Kind | Purpose |
|------|------|---------|
| `github-token-secret.yaml` | Secret | GitHub token for cloning and PR creation |
| `credentials-secret.yaml` | Secret | Claude OAuth token for the agents |
| `workspace.yaml` | Workspace | Git repository to clone |
| `task-a-implement.yaml` | Task | Step 1: implement a feature and open a PR |
| `task-b-review.yaml` | Task | Step 2: review the PR (depends on Task A) |

## Steps

1. **Edit the secrets** — replace placeholders in both `github-token-secret.yaml`
   and `credentials-secret.yaml` with your real tokens.

2. **Edit `workspace.yaml`** — set your repository URL and branch.

3. **Apply all resources at once:**

```bash
kubectl apply -f examples/06-task-dependencies/
```

4. **Watch the Tasks:**

```bash
kubectl get tasks -w
```

You should see Task A start immediately while Task B stays in `Waiting`.
Once Task A succeeds, Task B transitions to `Pending` and then `Running`.

5. **Stream logs for each task:**

```bash
# In one terminal:
axon logs implement-feature -f

# In another terminal (once Task B starts):
axon logs review-feature -f
```

6. **Cleanup:**

```bash
kubectl delete -f examples/06-task-dependencies/
```

## Key Behaviors

- **Automatic waiting** — Task B never starts until Task A succeeds. No
  polling or manual coordination needed.
- **Failure propagation** — if Task A fails, Task B immediately fails with
  a message indicating the dependency failure.
- **Cycle detection** — the controller detects circular dependencies at
  creation time and fails the task immediately.
- **Branch locking** — both tasks use the same branch. The controller
  ensures only one runs at a time (Task B waits for the lock too).
- **Output forwarding** — Task A's captured outputs (branch, PR URL) are
  injected into Task B's prompt via Go templates.

## Extending This Pattern

You can chain more than two tasks by adding additional `dependsOn`
references. For example, a three-step pipeline:

```yaml
# task-c-merge.yaml
apiVersion: axon.io/v1alpha1
kind: Task
metadata:
  name: merge-feature
spec:
  type: claude-code
  prompt: |
    The review for PR {{ index .Deps "review-feature" "Results" "pr" }}
    has been completed. If the review approved the changes, merge the PR.
  dependsOn:
    - review-feature
  credentials:
    type: oauth
    secretRef:
      name: claude-oauth-token
  workspaceRef:
    name: my-workspace
```

A task can also depend on multiple upstream tasks:

```yaml
dependsOn:
  - implement-backend
  - implement-frontend
```

In this case, the task waits until **all** listed dependencies succeed.
