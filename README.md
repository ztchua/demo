# Expense Tracker

A simple web application for tracking expenses with a Go backend and vanilla JavaScript frontend.

## Features

- Create, read, update, and delete expenses
- Categorize expenses
- SQLite database for local persistence
- Simple web UI

## Prerequisites

- Go 1.21 or later
- CGO enabled (required for go-sqlite3)

## Setup

1. Clone or navigate to the project directory:
   ```bash
   cd /Users/ztchua/dev/projects/demo
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run the server:
   ```bash
   go run main.go
   ```

4. Open your browser and navigate to:
   ```
   http://localhost:4000
   ```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/expenses` | Get all expenses |
| POST | `/expenses` | Create a new expense |
| GET | `/expenses/{id}` | Get a specific expense |
| PUT | `/expenses/{id}` | Update an expense |
| DELETE | `/expenses/{id}` | Delete an expense |

## Configuration

The server runs on port 4000 by default. To use a different port:

```bash
PORT=8080 go run main.go
```

## Database

The SQLite database file (`expenses.db`) is created automatically in the project directory on first run.
