"""Utility MCP Server - JSON CRUD operations inspired by main.go patterns."""

import json
import uuid
from datetime import datetime
from fastmcp import FastMCP

# Create the MCP server
mcp = FastMCP("utility-mcp-server")


# ============================================
# CREATE - Add new items to lists/objects
# ============================================

@mcp.tool()
def create_item(json_list: str, item: str, id_field: str = "id") -> str:
    """Add a new item to a JSON list, auto-generating an ID if needed.

    Inspired by createExpense() from main.go - inserts new records with timestamps.

    Args:
        json_list: JSON array string
        item: JSON object to add
        id_field: Field name for ID (auto-generated if missing)

    Returns:
        Updated JSON list with the new item
    """
    items = json.loads(json_list)
    new_item = json.loads(item)

    # Auto-generate ID if not present
    if id_field not in new_item or not new_item[id_field]:
        new_item[id_field] = str(uuid.uuid4())

    # Add timestamps (like main.go)
    new_item["created_at"] = datetime.utcnow().isoformat() + "Z"
    new_item["updated_at"] = datetime.utcnow().isoformat() + "Z"

    items.append(new_item)
    return json.dumps(items, indent=2)


@mcp.tool()
def batch_create(json_list: str, items: str, id_field: str = "id") -> str:
    """Add multiple items to a JSON list at once.

    Args:
        json_list: JSON array string
        items: JSON array of objects to add
        id_field: Field name for ID

    Returns:
        Updated JSON list with new items
    """
    list_items = json.loads(json_list)
    new_items = json.loads(items)

    now = datetime.utcnow().isoformat() + "Z"
    for item in new_items:
        if id_field not in item or not item[id_field]:
            item[id_field] = str(uuid.uuid4())
        item["created_at"] = now
        item["updated_at"] = now
        list_items.append(item)

    return json.dumps(list_items, indent=2)


# ============================================
# READ - Retrieve items from lists
# ============================================

@mcp.tool()
def get_all(json_list: str) -> str:
    """Get all items from a JSON list.

    Inspired by getExpenses() from main.go - returns full list ordered.

    Args:
        json_list: JSON array string

    Returns:
        The full JSON list
    """
    items = json.loads(json_list)
    return json.dumps(items, indent=2)


@mcp.tool()
def get_by_id(json_list: str, id_value: str, id_field: str = "id") -> str:
    """Get a single item by its ID field.

    Inspired by getExpense(id) from main.go - fetches one record.

    Args:
        json_list: JSON array string
        id_value: The ID value to find
        id_field: Field name to search (default: "id")

    Returns:
        The matching item as JSON, or error if not found
    """
    items = json.loads(json_list)
    for item in items:
        if str(item.get(id_field)) == str(id_value):
            return json.dumps(item, indent=2)
    return json.dumps({"error": f"Item with {id_field}='{id_value}' not found"})


@mcp.tool()
def filter_items(json_list: str, field: str, value: str) -> str:
    """Filter items where a field equals a value.

    Args:
        json_list: JSON array string
        field: Field name to check
        value: Value to match

    Returns:
        JSON array of matching items
    """
    items = json.loads(json_list)
    filtered = [item for item in items if str(item.get(field)) == str(value)]
    return json.dumps(filtered, indent=2)


@mcp.tool()
def search_items(json_list: str, query: str) -> str:
    """Search for items containing text in any string field.

    Args:
        json_list: JSON array string
        query: Text to search for

    Returns:
        JSON array of matching items
    """
    items = json.loads(json_list)
    query_lower = query.lower()
    matches = [
        item for item in items
        if any(
            isinstance(v, str) and query_lower in v.lower()
            for v in item.values()
        )
    ]
    return json.dumps(matches, indent=2)


@mcp.tool()
def sort_items(json_list: str, sort_field: str, descending: bool = False) -> str:
    """Sort items by a field.

    Inspired by getExpenses() ORDER BY clause in main.go.

    Args:
        json_list: JSON array string
        sort_field: Field to sort by
        descending: Sort descending if true

    Returns:
        Sorted JSON array
    """
    items = json.loads(json_list)

    def get_sort_key(item):
        value = item.get(sort_field, "")
        # Handle None values
        if value is None:
            return ""
        return value

    sorted_items = sorted(items, key=get_sort_key, reverse=descending)
    return json.dumps(sorted_items, indent=2)


# ============================================
# UPDATE - Modify existing items
# ============================================

@mcp.tool()
def update_item(json_list: str, id_value: str, updates: str, id_field: str = "id") -> str:
    """Update an item by ID, merging in new field values.

    Inspired by updateExpense(id, body) from main.go - updates record and refreshes updated_at.

    Args:
        json_list: JSON array string
        id_value: ID of item to update
        updates: JSON object with fields to update
        id_field: Field name for ID

    Returns:
        Updated JSON list, or error if not found
    """
    items = json.loads(json_list)
    update_data = json.loads(updates)

    for i, item in enumerate(items):
        if str(item.get(id_field)) == str(id_value):
            # Merge updates (like main.go does)
            items[i].update(update_data)
            # Always update the timestamp
            items[i]["updated_at"] = datetime.utcnow().isoformat() + "Z"
            # Keep original created_at
            if "created_at" not in items[i]:
                items[i]["created_at"] = items[i]["updated_at"]
            return json.dumps(items, indent=2)

    return json.dumps({"error": f"Item with {id_field}='{id_value}' not found"})


@mcp.tool()
def batch_update(json_list: str, ids: str, updates: str, id_field: str = "id") -> str:
    """Update multiple items by their IDs.

    Args:
        json_list: JSON array string
        ids: JSON array of ID values to update
        updates: JSON object with fields to update
        id_field: Field name for ID

    Returns:
        Updated JSON list
    """
    items = json.loads(json_list)
    target_ids = set(json.loads(ids))
    update_data = json.loads(updates)
    now = datetime.utcnow().isoformat() + "Z"

    updated_count = 0
    for i, item in enumerate(items):
        if str(item.get(id_field)) in target_ids:
            items[i].update(update_data)
            items[i]["updated_at"] = now
            updated_count += 1

    return json.dumps({
        "items": items,
        "updated_count": updated_count
    }, indent=2)


@mcp.tool()
def upsert_item(json_list: str, item: str, id_field: str = "id") -> str:
    """Update an item if it exists, or insert if it doesn't.

    Args:
        json_list: JSON array string
        item: JSON object (must contain id_field)
        id_field: Field name for ID

    Returns:
        Updated JSON list
    """
    items = json.loads(json_list)
    new_item = json.loads(item)
    item_id = new_item.get(id_field)

    if not item_id:
        # No ID provided - create new
        return create_item(json.dumps(items), json.dumps(new_item), id_field)

    # Try to update existing
    for i, item in enumerate(items):
        if str(item.get(id_field)) == str(item_id):
            items[i].update(new_item)
            items[i]["updated_at"] = datetime.utcnow().isoformat() + "Z"
            if "created_at" not in items[i]:
                items[i]["created_at"] = items[i]["updated_at"]
            return json.dumps(items, indent=2)

    # Not found - insert new
    return create_item(json.dumps(items), json.dumps(new_item), id_field)


# ============================================
# DELETE - Remove items
# ============================================

@mcp.tool()
def delete_item(json_list: str, id_value: str, id_field: str = "id") -> str:
    """Delete an item by ID.

    Inspired by deleteExpense(id) from main.go - removes record and returns 204/404 equivalent.

    Args:
        json_list: JSON array string
        id_value: ID of item to delete
        id_field: Field name for ID

    Returns:
        Updated JSON list, or error if not found
    """
    items = json.loads(json_list)
    original_length = len(items)

    items = [item for item in items if str(item.get(id_field)) != str(id_value)]

    if len(items) == original_length:
        return json.dumps({"error": f"Item with {id_field}='{id_value}' not found"})

    return json.dumps(items, indent=2)


@mcp.tool()
def batch_delete(json_list: str, ids: str, id_field: str = "id") -> str:
    """Delete multiple items by their IDs.

    Args:
        json_list: JSON array string
        ids: JSON array of ID values to delete
        id_field: Field name for ID

    Returns:
        Updated JSON list with delete count
    """
    items = json.loads(json_list)
    target_ids = set(json.loads(ids))
    original_length = len(items)

    items = [item for item in items if str(item.get(id_field)) not in target_ids]

    deleted_count = original_length - len(items)
    return json.dumps({
        "items": items,
        "deleted_count": deleted_count
    }, indent=2)


@mcp.tool()
def filter_delete(json_list: str, field: str, value: str) -> str:
    """Delete all items where a field equals a value.

    Args:
        json_list: JSON array string
        field: Field name to check
        value: Value to match for deletion

    Returns:
        Updated JSON list with delete count
    """
    items = json.loads(json_list)
    original_length = len(items)

    items = [item for item in items if str(item.get(field)) != str(value)]

    deleted_count = original_length - len(items)
    return json.dumps({
        "items": items,
        "deleted_count": deleted_count
    }, indent=2)


# ============================================
# VALIDATION - Like main.go field checks
# ============================================

@mcp.tool()
def validate_required(item: str, required_fields: list[str]) -> str:
    """Check that an item has all required fields.

    Inspired by main.go validation: `if e.Description == "" || e.Amount == 0`

    Args:
        item: JSON object to validate
        required_fields: List of field names that must be present and non-empty

    Returns:
        JSON with "valid" boolean and "missing" list
    """
    data = json.loads(item)
    missing = [
        field for field in required_fields
        if field not in data or not data[field]
    ]

    return json.dumps({
        "valid": len(missing) == 0,
        "missing": missing
    }, indent=2)


# ============================================
# Entry Point
# ============================================

def main():
    """Entry point for running the server."""
    mcp.run()


if __name__ == "__main__":
    main()
