#!/bin/bash
# axon_entrypoint.sh — Axon agent image interface implementation for
# Cursor CLI.
#
# Interface contract:
#   - First argument ($1): the task prompt
#   - CURSOR_API_KEY env var: API key for authentication
#   - AXON_MODEL env var: model name (optional)
#   - UID 61100: shared between git-clone init container and agent
#   - Working directory: /workspace/repo when a workspace is configured

set -uo pipefail

PROMPT="${1:?Prompt argument is required}"

ARGS=(
  "-p"
  "--force"
  "--trust"
  "--sandbox" "disabled"
  "--output-format" "stream-json"
  "$PROMPT"
)

if [ -n "${AXON_MODEL:-}" ]; then
  ARGS=("--model" "$AXON_MODEL" "${ARGS[@]}")
fi

# Write user-level instructions to both user config and workspace root.
# Cursor CLI may read AGENTS.md from the working directory.
if [ -n "${AXON_AGENTS_MD:-}" ]; then
  mkdir -p ~/.cursor
  printf '%s' "$AXON_AGENTS_MD" >~/.cursor/AGENTS.md
  printf '%s' "$AXON_AGENTS_MD" >/workspace/AGENTS.md
fi

# Install each plugin's skills into Cursor's .cursor/skills/ directory
# in the workspace so the CLI discovers them at runtime.
if [ -n "${AXON_PLUGIN_DIR:-}" ] && [ -d "${AXON_PLUGIN_DIR}" ]; then
  for plugindir in "${AXON_PLUGIN_DIR}"/*/; do
    [ -d "$plugindir" ] || continue
    if [ -d "${plugindir}skills" ]; then
      for skilldir in "${plugindir}skills"/*/; do
        [ -d "$skilldir" ] || continue
        skillname=$(basename "$skilldir")
        pluginname=$(basename "$plugindir")
        targetdir="/workspace/.cursor/skills/${pluginname}-${skillname}"
        mkdir -p "$targetdir"
        if [ -f "${skilldir}SKILL.md" ]; then
          cp "${skilldir}SKILL.md" "$targetdir/SKILL.md"
        fi
      done
    fi
  done
fi

# Write MCP server configuration to user-scoped ~/.cursor/mcp.json.
# The AXON_MCP_SERVERS JSON format matches Cursor's native format directly.
if [ -n "${AXON_MCP_SERVERS:-}" ]; then
  mkdir -p ~/.cursor
  node -e '
const fs = require("fs");
const cfgPath = require("os").homedir() + "/.cursor/mcp.json";
let existing = {};
try { existing = JSON.parse(fs.readFileSync(cfgPath, "utf8")); } catch {}
const mcp = JSON.parse(process.env.AXON_MCP_SERVERS);
existing.mcpServers = Object.assign(existing.mcpServers || {}, mcp.mcpServers || {});
fs.writeFileSync(cfgPath, JSON.stringify(existing, null, 2));
'
fi

agent "${ARGS[@]}" | tee /tmp/agent-output.jsonl
AGENT_EXIT_CODE=${PIPESTATUS[0]}

/axon/axon-capture

exit $AGENT_EXIT_CODE
