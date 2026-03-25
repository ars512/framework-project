package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"shop/handlers"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/products", handlers.GetProducts).Methods("GET")
	r.HandleFunc("/products", handlers.CreateProduct).Methods("POST")
	r.HandleFunc("/products/{id}", handlers.GetProduct).Methods("GET")
	r.HandleFunc("/products/{id}", handlers.UpdateProduct).Methods("PUT")
	r.HandleFunc("/products/{id}", handlers.DeleteProduct).Methods("DELETE")

	r.HandleFunc("/brands", handlers.GetBrands).Methods("GET")
	r.HandleFunc("/brands", handlers.CreateBrand).Methods("POST")

	r.HandleFunc("/categories", handlers.GetCategories).Methods("GET")
	r.HandleFunc("/categories", handlers.CreateCategory).Methods("POST")

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}
