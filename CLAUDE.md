# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Important: Do NOT Delete Files

**NEVER use `rm`, `rmdir`, or similar commands to delete files without explicit user permission.** Always ask the user before deleting any files in this project.

## Commands

```bash
# Run the server
go run main.go

# Build the binary
go build main.go

# Run the built binary
./main

# Run tests
go test -v ./...

# Run tests with coverage
go test -cover ./...
```

The server runs on port 4000 by default (configurable via `PORT` env var).

## Architecture

This is an expense tracking web application with a Go backend and vanilla JavaScript frontend.

**Backend (`main.go`):** Single-file REST API with SQLite persistence

**Key design decisions:**
- `Server` struct holds the database connection for dependency injection (enables testing)
- All handlers are methods on `Server` rather than standalone functions
- Database file (`expenses.db`) is created in the working directory on first run
- Dates are stored as RFC3339 strings in SQLite (SQLite lacks native DATETIME type)
- All timestamp fields (`date`, `created_at`, `updated_at`) are parsed from string format after SQL queries
- ID extraction from URL path uses manual string slicing - `/expenses/` prefix is stripped to get the ID

**Handler routing:**
- `/expenses` - Collection endpoint (GET list, POST create)
- `/expenses/` - Individual resource endpoint (GET, PUT, DELETE by ID)

The route pattern relies on Go's `http.HandleFunc` behavior where `/expenses/` matches any path with that prefix. The ID is parsed from the tail of the path string.

**Dependencies:**
- `github.com/mattn/go-sqlite3` - CGO-based SQLite driver (requires CGO enabled)

**Frontend (`static/`):**
- `index.html` - Single-page UI with form and expense table
- `app.js` - Vanilla JavaScript for CRUD operations, no build step required

**Testing:**
- `main_test.go` - Test helper functions (SetupTestDB, CleanupTestDB, SeedExpenses, etc.)
- `handlers_test.go` - Comprehensive endpoint tests (23 test cases, ~77% coverage)
- Tests use in-memory SQLite (`:memory:`) for isolation

**Project structure:**
```
├── main.go              # Backend server with API handlers
├── main_test.go         # Test helpers and setup
├── handlers_test.go     # Endpoint tests
├── static/
│   ├── index.html      # Frontend UI
│   └── app.js          # Frontend logic
├── expenses.db         # SQLite database (created on first run)
├── go.mod              # Go module definition
└── .claude/
    └── settings.json   # Project settings (PostToolUse hook runs tests)
```

**CORS:** API endpoints have CORS headers enabled via `withCORS()` middleware to allow frontend communication.
