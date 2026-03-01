package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// Bad Example: Everything in one place, no layers
// Handler, business logic, and DB queries are all mixed together

type Customer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var db *sql.DB

func init() {
	// ❌ DB connection in init — hard to test, can't return error
	var err error
	db, err = sql.Open("mysql", "user:pass@/shop")
	if err != nil {
		panic(err)
	}
}

// ❌ HTTP handler, business logic, and DB query all mixed in one function
func createCustomerHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// ❌ Business logic mixed in handler
	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	if len(req.Name) > 50 {
		http.Error(w, "name too long", http.StatusBadRequest)
		return
	}

	// ❌ DB query mixed in handler
	result, err := db.Exec("INSERT INTO customers (name) VALUES (?)", req.Name)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	id, _ := result.LastInsertId()
	_ = id

	// Return response
	customer := Customer{Name: req.Name}
	json.NewEncoder(w).Encode(customer)
}

func getCustomerHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	// ❌ Business logic mixed in handler
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	// ❌ DB query mixed in handler
	var customer Customer
	err := db.QueryRow("SELECT id, name FROM customers WHERE id = ?", id).
		Scan(&customer.ID, &customer.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(customer)
}

func main() {
	http.HandleFunc("/customers", createCustomerHandler)
	http.HandleFunc("/customers/get", getCustomerHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Problems:
// 1. Cannot test business logic without a real DB
// 2. Cannot swap DB implementation
// 3. Business rules scattered across handlers
// 4. init() makes testing even harder
// 5. No separation of concerns — changes in DB schema touch handler code
