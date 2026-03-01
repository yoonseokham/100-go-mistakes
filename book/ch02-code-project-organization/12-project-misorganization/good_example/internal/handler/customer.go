package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/yourusername/100-go-mistakes/book/ch02-code-project-organization/12-project-misorganization/good_example/internal/service"
)

// CustomerHandler knows only about HTTP and Service interface
// It does NOT know about DB or business rules

// Defined on consumer side — Chapter 6/7 principle
type customerService interface {
	CreateCustomer(name string) (service.Customer, error)
	GetCustomer(id string) (service.Customer, error)
}

type CustomerHandler struct {
	svc customerService
}

func NewCustomerHandler(svc customerService) *CustomerHandler {
	return &CustomerHandler{svc: svc}
}

func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Handler's only job: parse request, call service, write response
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	customer, err := h.svc.CreateCustomer(req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

func (h *CustomerHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	customer, err := h.svc.GetCustomer(id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}
