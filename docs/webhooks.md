# GitHub Webhook Support

Kelos supports GitHub webhooks for work item discovery. Instead of polling the GitHub API, Kelos can receive push notifications when issues or pull requests are created or updated.

## Architecture

The webhook system consists of three components:

### 1. WebhookEvent CRD

Webhook payloads are stored as `WebhookEvent` custom resources in Kubernetes, providing:

- **Persistence**: Events survive pod restarts (stored in etcd)
- **Auditability**: All events are visible via `kubectl get webhookevents`
- **Processing tracking**: Events are marked as processed after discovery

### 2. Webhook Receiver (kelos-webhook-receiver)

An HTTP server that:
- Listens on `/webhook/github`
- Validates GitHub webhook signatures (HMAC-SHA256)
- Creates `WebhookEvent` CRD instances
- Returns 202 Accepted

Deploy as a Deployment with a LoadBalancer Service to expose it publicly.

### 3. GitHubWebhookSource

The `GitHubWebhookSource` implementation:
- Lists unprocessed `WebhookEvent` resources
- Parses GitHub webhook payloads into `WorkItem` format
- Applies filters (labels, state, etc.)
- Marks events as processed

## GitHub Webhook Setup

### 1. Deploy the webhook receiver

```bash
kubectl apply -f examples/taskspawner-github-webhook.yaml
```

This creates:
- `kelos-webhook-receiver` Deployment
- LoadBalancer Service
- RBAC for creating WebhookEvent resources

### 2. Get the external URL

```bash
kubectl get service kelos-webhook-receiver -o jsonpath='{.status.loadBalancer.ingress[0].hostname}'
```

### 3. Configure GitHub webhook

In your GitHub repository settings:
1. Go to Settings → Webhooks → Add webhook
2. **Payload URL**: `http://<your-loadbalancer>/webhook/github`
3. **Content type**: `application/json`
4. **Secret**: Set a secret and store it in the `github-webhook-secret` Secret
5. **Events**: Select "Issues" and "Pull requests"
6. Click "Add webhook"

### 4. Create a TaskSpawner with githubWebhook

```yaml
apiVersion: kelos.dev/v1alpha1
kind: TaskSpawner
metadata:
  name: my-webhook-spawner
spec:
  when:
    githubWebhook:
      namespace: default
      labels:
        - "kelos-task"
  taskTemplate:
    type: claude-code
    credentials:
      type: api-key
      secretRef:
        name: anthropic-api-key
    workspaceRef:
      name: my-workspace
    promptTemplate: |
      {{ .Title }}
      {{ .Body }}
```

## Webhook Signature Validation

For GitHub webhooks, the receiver validates the `X-Hub-Signature-256` header using HMAC-SHA256.

Set the `GITHUB_WEBHOOK_SECRET` environment variable on the webhook receiver to enable validation:

```yaml
env:
- name: GITHUB_WEBHOOK_SECRET
  valueFrom:
    secretKeyRef:
      name: github-webhook-secret
      key: secret
```

If the secret is not set, signature validation is skipped (development mode only).
