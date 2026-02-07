#!/bin/bash
# axon_entrypoint.sh - Standard axon agent image entrypoint for Claude Code.
#
# This script implements the axon agent image interface:
#   - Receives the prompt as the first argument ($1)
#   - Uses AXON_MODEL environment variable for the model (if set)
#
# Reserved environment variables set by axon:
#   AXON_PROMPT - the task prompt (same as $1)
#   AXON_MODEL  - the model to use (may be empty)
set -euo pipefail

PROMPT="${1:?Prompt argument is required}"

ARGS=(
    "--dangerously-skip-permissions"
    "--output-format" "stream-json"
    "--verbose"
    "-p" "${PROMPT}"
)

if [ -n "${AXON_MODEL:-}" ]; then
    ARGS+=("--model" "${AXON_MODEL}")
fi

exec claude "${ARGS[@]}"
