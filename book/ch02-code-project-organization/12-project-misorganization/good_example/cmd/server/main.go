package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/yourusername/100-go-mistakes/book/ch02-code-project-organization/12-project-misorganization/good_example/internal/handler"
	"github.com/yourusername/100-go-mistakes/book/ch02-code-project-organization/12-project-misorganization/good_example/internal/repository"
	"github.com/yourusername/100-go-mistakes/book/ch02-code-project-organization/12-project-misorganization/good_example/internal/service"
)

func main() {
	// Dependency wiring happens ONLY in main
	// Each layer receives what it needs via constructor (DI)

	db, err := sql.Open("mysql", "user:pass@/shop")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Layer 3: Repository — knows DB
	repo := repository.NewCustomerRepository(db)

	// Layer 2: Service — knows Repository interface
	svc := service.NewCustomerService(repo)

	// Layer 1: Handler — knows Service interface
	h := handler.NewCustomerHandler(svc)

	http.HandleFunc("/customers", h.Create)
	http.HandleFunc("/customers/get", h.Get)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Dependency flow:
//   main → Handler → Service → Repository → DB
//
// Each layer only knows the layer below via interface:
//   Handler  knows: customerService interface
//   Service  knows: CustomerRepo interface
//   Repository knows: *sql.DB
