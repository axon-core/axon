#!/bin/bash
# axon_entrypoint.sh — Axon agent image interface implementation for
# OpenAI Codex CLI.
#
# Interface contract:
#   - First argument ($1): the task prompt
#   - AXON_MODEL env var: model name (optional)
#   - UID 61100: shared between git-clone init container and agent
#   - Working directory: /workspace/repo when a workspace is configured

set -euo pipefail

PROMPT="${1:?Prompt argument is required}"

ARGS=(
    "exec"
    "--dangerously-bypass-approvals-and-sandbox"
    "--json"
    "$PROMPT"
)

if [ -n "${AXON_MODEL:-}" ]; then
    ARGS+=("--model" "$AXON_MODEL")
fi

exec codex "${ARGS[@]}"
