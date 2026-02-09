#!/bin/bash
# axon_entrypoint.sh — Axon agent image interface implementation for
# Google Gemini CLI.
#
# Interface contract:
#   - First argument ($1): the task prompt
#   - AXON_MODEL env var: model name (optional)
#   - UID 61100: shared between git-clone init container and agent
#   - Working directory: /workspace/repo when a workspace is configured

set -uo pipefail

PROMPT="${1:?Prompt argument is required}"

ARGS=(
    "--yolo"
    "--output-format" "stream-json"
    "-p" "$PROMPT"
)

if [ -n "${AXON_MODEL:-}" ]; then
    ARGS+=("--model" "$AXON_MODEL")
fi

gemini "${ARGS[@]}"
AGENT_EXIT_CODE=$?

/axon/capture-outputs.sh

exit $AGENT_EXIT_CODE
