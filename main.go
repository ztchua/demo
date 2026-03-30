package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Expense represents an expense record
type Expense struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Server holds the application state
type Server struct {
	DB *sql.DB
}

// withCORS wraps a handler to add CORS headers
func withCORS(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		fn(w, r)
	}
}

func main() {
	db, err := sql.Open("sqlite3", "./expenses.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := createTable(db); err != nil {
		log.Fatal(err)
	}

	server := &Server{DB: db}

	// API handlers with CORS
	http.HandleFunc("/expenses", withCORS(server.expensesHandler))
	http.HandleFunc("/expenses/", withCORS(server.expenseByIDHandler))

	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	port := "4000"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	log.Printf("Server starting on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func createTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS expenses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		description TEXT NOT NULL,
		amount REAL NOT NULL,
		category TEXT,
		date TEXT NOT NULL,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL
	);`

	_, err := db.Exec(query)
	return err
}

func (s *Server) expensesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getExpenses(w, r)
	case http.MethodPost:
		s.createExpense(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) expenseByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/expenses/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid expense ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.getExpense(w, r, id)
	case http.MethodPut:
		s.updateExpense(w, r, id)
	case http.MethodDelete:
		s.deleteExpense(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) getExpenses(w http.ResponseWriter, r *http.Request) {
	rows, err := s.DB.Query("SELECT id, description, amount, category, date, created_at, updated_at FROM expenses ORDER BY date DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	expenses := []Expense{}
	for rows.Next() {
		var e Expense
		var dateStr, createdAtStr, updatedAtStr string
		err := rows.Scan(&e.ID, &e.Description, &e.Amount, &e.Category, &dateStr, &createdAtStr, &updatedAtStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		e.Date, _ = time.Parse(time.RFC3339, dateStr)
		e.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		e.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAtStr)
		expenses = append(expenses, e)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}

func (s *Server) getExpense(w http.ResponseWriter, r *http.Request, id int) {
	var e Expense
	var dateStr, createdAtStr, updatedAtStr string

	err := s.DB.QueryRow(
		"SELECT id, description, amount, category, date, created_at, updated_at FROM expenses WHERE id = ?",
		id,
	).Scan(&e.ID, &e.Description, &e.Amount, &e.Category, &dateStr, &createdAtStr, &updatedAtStr)

	if err == sql.ErrNoRows {
		http.Error(w, "Expense not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	e.Date, _ = time.Parse(time.RFC3339, dateStr)
	e.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
	e.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAtStr)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(e)
}

func (s *Server) createExpense(w http.ResponseWriter, r *http.Request) {
	var e Expense
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if e.Description == "" || e.Amount == 0 {
		http.Error(w, "Description and amount are required", http.StatusBadRequest)
		return
	}

	now := time.Now()
	if e.Date.IsZero() {
		e.Date = now
	}
	e.CreatedAt = now
	e.UpdatedAt = now

	result, err := s.DB.Exec(
		"INSERT INTO expenses (description, amount, category, date, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		e.Description, e.Amount, e.Category, e.Date.Format(time.RFC3339), e.CreatedAt.Format(time.RFC3339), e.UpdatedAt.Format(time.RFC3339),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	e.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(e)
}

func (s *Server) updateExpense(w http.ResponseWriter, r *http.Request, id int) {
	var e Expense
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	e.ID = id
	e.UpdatedAt = time.Now()

	result, err := s.DB.Exec(
		"UPDATE expenses SET description = ?, amount = ?, category = ?, date = ?, updated_at = ? WHERE id = ?",
		e.Description, e.Amount, e.Category, e.Date.Format(time.RFC3339), e.UpdatedAt.Format(time.RFC3339), id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Expense not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(e)
}

func (s *Server) deleteExpense(w http.ResponseWriter, r *http.Request, id int) {
	result, err := s.DB.Exec("DELETE FROM expenses WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Expense not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
