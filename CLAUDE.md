# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run the server
go run main.go

# Build the binary
go build main.go

# Run the built binary
./main
```

The server runs on port 4000 by default (configurable via `PORT` env var).

## Architecture

This is an expense tracking web application with a Go backend and vanilla JavaScript frontend.

**Backend (`main.go`):** Single-file REST API with SQLite persistence

**Key design decisions:**
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

**Project structure:**
```
├── main.go           # Backend server with API handlers
├── static/
│   ├── index.html   # Frontend UI
│   └── app.js       # Frontend logic
├── expenses.db      # SQLite database (created on first run)
└── go.mod           # Go module definition
```

**CORS:** API endpoints have CORS headers enabled via `withCORS()` middleware to allow frontend communication.
