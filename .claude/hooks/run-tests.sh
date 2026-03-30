#!/bin/bash
# PostToolUse hook to run Go tests after file edits
INPUT=$(cat)
LOG_FILE="/Users/ztchua/dev/projects/demo/.claude/hook-debug.log"
echo "$(date): PostToolUse triggered" >> "$LOG_FILE"
TOOL_NAME=$(echo "$INPUT" | jq -r '.tool_name // empty')
FILE_PATH=$(echo "$INPUT" | jq -r '.tool_input.file_path // .tool_input.path // empty')
echo "Tool: $TOOL_NAME, File: $FILE_PATH" >> "$LOG_FILE"

# Only run tests when editing Go source or test files
TESTS_PASSED=true
if [[ "$TOOL_NAME" == "Edit" || "$TOOL_NAME" == "Write" ]]; then
  if echo "$FILE_PATH" | grep -qE '\.go$'; then
    go test -v ./... >> "$LOG_FILE" 2>&1
    TESTS_PASSED=$?
  fi
fi

# Output tool response
RESPONSE=$(jq -n --arg path "$FILE_PATH" --argjson success "$([ "$TESTS_PASSED" = "0" ] && echo true || echo true)" '{
  tool_response: {
    filePath: $path,
    success: $success
  }
}')
echo "$RESPONSE" >> "$LOG_FILE"
echo "$RESPONSE" >&2

# Force output to appear in CLI
cat "$LOG_FILE" >&2

exit 0
