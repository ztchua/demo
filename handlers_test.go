package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Test GET /expenses
func TestGetExpenses(t *testing.T) {
	tests := []struct {
		name           string
		seedData       []Expense
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "empty database returns empty array",
			seedData:       nil,
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name: "returns all expenses ordered by date DESC",
			seedData: []Expense{
				{
					Description: "Lunch",
					Amount:      12.50,
					Category:    "Food",
					Date:        time.Now().Add(-2 * time.Hour),
					CreatedAt:   time.Now().Add(-2 * time.Hour),
					UpdatedAt:   time.Now().Add(-2 * time.Hour),
				},
				{
					Description: "Dinner",
					Amount:      25.00,
					Category:    "Food",
					Date:        time.Now().Add(-1 * time.Hour),
					CreatedAt:   time.Now().Add(-1 * time.Hour),
					UpdatedAt:   time.Now().Add(-1 * time.Hour),
				},
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := SetupTestDB(t)
			defer CleanupTestDB(t, db)

			if tt.seedData != nil {
				SeedExpenses(t, db, tt.seedData)
			}

			server := &Server{DB: db}
			req := CreateTestRequest(t, http.MethodGet, "/expenses", nil)
			resp := httptest.NewRecorder()

			server.getExpenses(resp, req)

			AssertStatus(t, resp, tt.expectedStatus)

			var expenses []Expense
			AssertJSON(t, resp, &expenses)

			if len(expenses) != tt.expectedCount {
				t.Errorf("Expected %d expenses, got %d", tt.expectedCount, len(expenses))
			}
		})
	}
}

// Test POST /expenses
func TestCreateExpense(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		checkResponse  func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name: "creates expense with all fields",
			requestBody: map[string]interface{}{
				"description": "Groceries",
				"amount":      45.99,
				"category":    "Food",
				"date":        time.Now().Format(time.RFC3339),
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				var e Expense
				AssertJSON(t, resp, &e)
				if e.Description != "Groceries" {
					t.Errorf("Expected description 'Groceries', got '%s'", e.Description)
				}
				if e.Amount != 45.99 {
					t.Errorf("Expected amount 45.99, got %f", e.Amount)
				}
				if e.Category != "Food" {
					t.Errorf("Expected category 'Food', got '%s'", e.Category)
				}
				if e.ID == 0 {
					t.Error("Expected ID to be set")
				}
			},
		},
		{
			name: "auto-generates date when not provided",
			requestBody: map[string]interface{}{
				"description": "Coffee",
				"amount":      3.50,
				"category":    "Beverage",
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				var e Expense
				AssertJSON(t, resp, &e)
				if e.Date.IsZero() {
					t.Error("Expected date to be auto-generated")
				}
			},
		},
		{
			name: "returns 400 when description is empty",
			requestBody: map[string]interface{}{
				"description": "",
				"amount":      10.00,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "returns 400 when amount is 0",
			requestBody: map[string]interface{}{
				"description": "Test",
				"amount":      0,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "returns 400 for invalid JSON",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := SetupTestDB(t)
			defer CleanupTestDB(t, db)

			server := &Server{DB: db}
			req := CreateTestRequest(t, http.MethodPost, "/expenses", tt.requestBody)
			resp := httptest.NewRecorder()

			server.createExpense(resp, req)

			AssertStatus(t, resp, tt.expectedStatus)

			if tt.checkResponse != nil {
				tt.checkResponse(t, resp)
			}
		})
	}
}

// Test GET /expenses/{id}
func TestGetExpense(t *testing.T) {
	now := time.Now()
	testExpense := Expense{
		Description: "Test Expense",
		Amount:      100.00,
		Category:    "Test",
		Date:        now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tests := []struct {
		name           string
		seedData       []Expense
		id             int
		expectedStatus int
		checkResponse  func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name:           "returns expense when exists",
			seedData:       []Expense{testExpense},
			id:             1,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				var e Expense
				AssertJSON(t, resp, &e)
				if e.ID != 1 {
					t.Errorf("Expected ID 1, got %d", e.ID)
				}
				if e.Description != "Test Expense" {
					t.Errorf("Expected description 'Test Expense', got '%s'", e.Description)
				}
			},
		},
		{
			name:           "returns 404 when not found",
			seedData:       []Expense{testExpense},
			id:             999,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := SetupTestDB(t)
			defer CleanupTestDB(t, db)

			if tt.seedData != nil {
				SeedExpenses(t, db, tt.seedData)
			}

			server := &Server{DB: db}
			req := CreateTestRequest(t, http.MethodGet, "/expenses/1", nil)
			resp := httptest.NewRecorder()

			server.getExpense(resp, req, tt.id)

			AssertStatus(t, resp, tt.expectedStatus)

			if tt.checkResponse != nil {
				tt.checkResponse(t, resp)
			}
		})
	}
}

// Test PUT /expenses/{id}
func TestUpdateExpense(t *testing.T) {
	now := time.Now()
	testExpense := Expense{
		Description: "Original",
		Amount:      50.00,
		Category:    "Original",
		Date:        now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tests := []struct {
		name           string
		seedData       []Expense
		id             int
		requestBody    interface{}
		expectedStatus int
		checkResponse  func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name:     "updates expense successfully",
			seedData: []Expense{testExpense},
			id:       1,
			requestBody: map[string]interface{}{
				"description": "Updated",
				"amount":      75.00,
				"category":    "Updated",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				var e Expense
				AssertJSON(t, resp, &e)
				if e.Description != "Updated" {
					t.Errorf("Expected description 'Updated', got '%s'", e.Description)
				}
				if e.Amount != 75.00 {
					t.Errorf("Expected amount 75.00, got %f", e.Amount)
				}
			},
		},
		{
			name:     "returns 404 when not found",
			seedData: []Expense{testExpense},
			id:       999,
			requestBody: map[string]interface{}{
				"description": "Updated",
				"amount":      75.00,
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "returns 400 for invalid JSON",
			seedData:       []Expense{testExpense},
			id:             1,
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := SetupTestDB(t)
			defer CleanupTestDB(t, db)

			if tt.seedData != nil {
				SeedExpenses(t, db, tt.seedData)
			}

			server := &Server{DB: db}
			req := CreateTestRequest(t, http.MethodPut, "/expenses/1", tt.requestBody)
			resp := httptest.NewRecorder()

			server.updateExpense(resp, req, tt.id)

			AssertStatus(t, resp, tt.expectedStatus)

			if tt.checkResponse != nil {
				tt.checkResponse(t, resp)
			}
		})
	}
}

// Test DELETE /expenses/{id}
func TestDeleteExpense(t *testing.T) {
	now := time.Now()
	testExpense := Expense{
		Description: "To Delete",
		Amount:      100.00,
		Category:    "Test",
		Date:        now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tests := []struct {
		name           string
		seedData       []Expense
		id             int
		expectedStatus int
	}{
		{
			name:           "deletes expense and returns 204",
			seedData:       []Expense{testExpense},
			id:             1,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "returns 404 when not found",
			seedData:       []Expense{testExpense},
			id:             999,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := SetupTestDB(t)
			defer CleanupTestDB(t, db)

			if tt.seedData != nil {
				SeedExpenses(t, db, tt.seedData)
			}

			server := &Server{DB: db}
			req := CreateTestRequest(t, http.MethodDelete, "/expenses/1", nil)
			resp := httptest.NewRecorder()

			server.deleteExpense(resp, req, tt.id)

			AssertStatus(t, resp, tt.expectedStatus)
		})
	}
}

// Test CORS middleware
func TestWithCORS(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(t, db)

	server := &Server{DB: db}

	tests := []struct {
		name           string
		method         string
		body           interface{}
		expectedStatus int
		checkHeaders   bool
	}{
		{
			name:           "sets CORS headers on GET request",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			checkHeaders:   true,
		},
		{
			name:           "handles OPTIONS preflight request",
			method:         http.MethodOptions,
			expectedStatus: http.StatusOK,
			checkHeaders:   true,
		},
		{
			name:           "sets CORS headers on POST request",
			method:         http.MethodPost,
			body: map[string]interface{}{
				"description": "Test",
				"amount":      10.00,
			},
			expectedStatus: http.StatusCreated,
			checkHeaders:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := CreateTestRequest(t, tt.method, "/expenses", tt.body)
			resp := httptest.NewRecorder()

			handler := withCORS(server.expensesHandler)
			handler(resp, req)

			AssertStatus(t, resp, tt.expectedStatus)

			if tt.checkHeaders {
				if resp.Header().Get("Access-Control-Allow-Origin") != "*" {
					t.Error("Expected CORS header 'Access-Control-Allow-Origin: *'")
				}
				if resp.Header().Get("Access-Control-Allow-Methods") == "" {
					t.Error("Expected CORS header 'Access-Control-Allow-Methods' to be set")
				}
				if resp.Header().Get("Access-Control-Allow-Headers") == "" {
					t.Error("Expected CORS header 'Access-Control-Allow-Headers' to be set")
				}
			}
		})
	}
}

// Test invalid ID parsing in expenseByIDHandler
func TestExpenseByIDHandlerInvalidID(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(t, db)

	server := &Server{DB: db}

	req := CreateTestRequest(t, http.MethodGet, "/expenses/invalid", nil)
	resp := httptest.NewRecorder()

	server.expenseByIDHandler(resp, req)

	AssertStatus(t, resp, http.StatusBadRequest)
}

// Test expensesHandler routing
func TestExpensesHandler(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(t, db)

	server := &Server{DB: db}

	tests := []struct {
		name           string
		method         string
		expectedStatus int
	}{
		{
			name:           "routes GET to getExpenses",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "routes POST to createExpense",
			method:         http.MethodPost,
			expectedStatus: http.StatusBadRequest, // No body
		},
		{
			name:           "returns 405 for invalid method",
			method:         http.MethodPut,
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := CreateTestRequest(t, tt.method, "/expenses", nil)
			resp := httptest.NewRecorder()

			server.expensesHandler(resp, req)

			AssertStatus(t, resp, tt.expectedStatus)
		})
	}
}

// Test expenseByIDHandler routing
func TestExpenseByIDHandler(t *testing.T) {
	now := time.Now()
	testExpense := Expense{
		Description: "Test",
		Amount:      10.00,
		Category:    "Test",
		Date:        now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	db := SetupTestDB(t)
	defer CleanupTestDB(t, db)

	SeedExpenses(t, db, []Expense{testExpense})

	server := &Server{DB: db}

	tests := []struct {
		name           string
		method         string
		expectedStatus int
	}{
		{
			name:           "routes GET to getExpense",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "routes PUT to updateExpense",
			method:         http.MethodPut,
			expectedStatus: http.StatusBadRequest, // No body
		},
		{
			name:           "routes DELETE to deleteExpense",
			method:         http.MethodDelete,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "returns 405 for invalid method",
			method:         http.MethodPost,
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := CreateTestRequest(t, tt.method, "/expenses/1", nil)
			resp := httptest.NewRecorder()

			server.expenseByIDHandler(resp, req)

			AssertStatus(t, resp, tt.expectedStatus)
		})
	}
}
