package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// SetupTestDB creates an in-memory SQLite database for testing
func SetupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	if err := createTable(db); err != nil {
		db.Close()
		t.Fatalf("Failed to create test table: %v", err)
	}

	return db
}

// CleanupTestDB closes the test database
func CleanupTestDB(t *testing.T, db *sql.DB) {
	t.Helper()
	if err := db.Close(); err != nil {
		t.Errorf("Failed to close test database: %v", err)
	}
}

// SeedExpenses inserts test expense records into the database
func SeedExpenses(t *testing.T, db *sql.DB, expenses []Expense) {
	t.Helper()
	for _, e := range expenses {
		_, err := db.Exec(
			"INSERT INTO expenses (description, amount, category, date, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
			e.Description, e.Amount, e.Category, e.Date.Format(time.RFC3339), e.CreatedAt.Format(time.RFC3339), e.UpdatedAt.Format(time.RFC3339),
		)
		if err != nil {
			t.Fatalf("Failed to seed expense: %v", err)
		}
	}
}

// AssertJSON asserts the response has correct content-type and decodes JSON
func AssertJSON(t *testing.T, resp *httptest.ResponseRecorder, v interface{}) {
	t.Helper()

	contentType := resp.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		t.Errorf("Failed to decode JSON response: %v", err)
	}
}

// AssertStatus asserts the HTTP status code
func AssertStatus(t *testing.T, resp *httptest.ResponseRecorder, expected int) {
	t.Helper()

	if resp.Code != expected {
		t.Errorf("Expected status %d, got %d. Body: %s", expected, resp.Code, resp.Body.String())
	}
}

// CreateTestRequest creates an HTTP request for testing
func CreateTestRequest(t *testing.T, method, url string, body interface{}) *http.Request {
	t.Helper()

	var reqBody *bytes.Buffer
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	} else {
		reqBody = &bytes.Buffer{}
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req
}
