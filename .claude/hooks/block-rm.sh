#!/bin/bash
# PreToolUse hook to block all file deletion commands
COMMAND=$(jq -r '.tool_input.command // empty')

if echo "$COMMAND" | grep -qE '\brm\b|\brmdir\b'; then
  jq -n --arg cmd "$COMMAND" '{
    hookSpecificOutput: {
      hookEventName: "PreToolUse",
      permissionDecision: "deny",
      permissionDecisionReason: ("Deletion command blocked by hook: " + $cmd)
    }
  }'
fi

exit 0
