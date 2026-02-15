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

# Write user-level instructions (global scope read by Gemini CLI)
if [ -n "${AXON_AGENTS_MD:-}" ]; then
  mkdir -p ~/.gemini
  printf '%s' "$AXON_AGENTS_MD" >~/.gemini/GEMINI.md
fi

# Install each plugin as a Gemini extension with skills and agents
if [ -n "${AXON_PLUGIN_DIR:-}" ] && [ -d "${AXON_PLUGIN_DIR}" ]; then
  for plugindir in "${AXON_PLUGIN_DIR}"/*/; do
    [ -d "$plugindir" ] || continue
    pluginname=$(basename "$plugindir")
    # Sanitize plugin name for safe JSON interpolation
    safename=$(printf '%s' "$pluginname" | tr -d '"\\\n\r')
    extdir="$HOME/.gemini/extensions/${pluginname}"
    mkdir -p "$extdir"
    printf '{"name":"%s"}' "$safename" >"$extdir/gemini-extension.json"
    # Copy skills directory
    if [ -d "${plugindir}skills" ]; then
      cp -r "${plugindir}skills" "$extdir/skills"
    fi
    # Copy agents directory
    if [ -d "${plugindir}agents" ]; then
      cp -r "${plugindir}agents" "$extdir/agents"
    fi
  done
fi

# Write MCP server configuration to Gemini settings.
# AXON_MCP_SERVERS contains JSON with an "mcpServers" key that Gemini
# settings.json accepts directly. Merge with existing settings if present.
if [ -n "${AXON_MCP_SERVERS:-}" ]; then
  settings_file="$HOME/.gemini/settings.json"
  if [ -f "$settings_file" ]; then
    # Merge mcpServers into existing settings using a small Node.js helper.
    # Read AXON_MCP_SERVERS from the environment to avoid exposing
    # potentially sensitive headers in process argument lists.
    node -e '
const fs = require("fs");
const existing = JSON.parse(fs.readFileSync(process.argv[1], "utf8"));
const mcp = JSON.parse(process.env.AXON_MCP_SERVERS);
existing.mcpServers = Object.assign(existing.mcpServers || {}, mcp.mcpServers || {});
fs.writeFileSync(process.argv[1], JSON.stringify(existing, null, 2));
' "$settings_file"
  else
    mkdir -p ~/.gemini
    printf '%s' "$AXON_MCP_SERVERS" >"$settings_file"
  fi
fi

gemini "${ARGS[@]}" | tee /tmp/agent-output.jsonl
AGENT_EXIT_CODE=${PIPESTATUS[0]}

/axon/axon-capture

exit $AGENT_EXIT_CODE
