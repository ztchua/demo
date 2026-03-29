# Utility MCP Server

A JSON CRUD utility server built with FastMCP. Tools inspired by standard REST patterns (create, read, update, delete).

## Installation

```bash
cd mcp_server
pip install -e .
```

## Running

```bash
# Direct (stdio mode)
python -m utility_mcp_server

# Or using the installed script
utility-mcp-server
```

## Configuration (Claude Code)

Add to your Claude Code settings (`~/.claude/settings.json`):

```json
{
  "mcpServers": {
    "utility": {
      "command": "uvx",
      "args": ["--directory", "/Users/ztchua/dev/projects/demo/mcp_server", "utility-mcp-server"]
    }
  }
}
```

Or with Python directly:

```json
{
  "mcpServers": {
    "utility": {
      "command": "python",
      "args": ["/Users/ztchua/dev/projects/demo/mcp_server/utility_mcp_server.py"]
    }
  }
}
```

## Tools

| Tool | Description |
|------|-------------|
| `create_item` | Add item to list, auto-generates ID + timestamps |
| `batch_create` | Add multiple items at once |
| `get_all` | Get all items |
| `get_by_id` | Get single item by ID |
| `filter_items` | Filter by field=value |
| `search_items` | Full-text search across string fields |
| `sort_items` | Sort by field |
| `update_item` | Update item by ID (refreshes updated_at) |
| `batch_update` | Update multiple items |
| `upsert_item` | Update if exists, insert if not |
| `delete_item` | Delete by ID |
| `batch_delete` | Delete multiple items |
| `filter_delete` | Delete all matching field=value |
| `validate_required` | Check required fields present |
