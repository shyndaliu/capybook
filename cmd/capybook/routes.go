package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() *mux.Router {
	r := mux.NewRouter()
	v1 := r.PathPrefix("/api/v1").Subrouter()
	v1.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)

	//Healthcheck
	v1.HandleFunc("/healthcheck", app.healthcheckHandler).Methods("GET")

	// Book Singleton
	// Create a new book
	v1.HandleFunc("/books", app.createBookHandler).Methods("POST")
	// Get all books
	//v1.HandleFunc("/books", app.getBookHandler).Methods("GET")
	//Get specific book
	v1.HandleFunc("/books/{id}", app.getBookHandler).Methods("GET")
	// Update a specific book
	v1.HandleFunc("/books/{id}", app.updateBookHandler).Methods("PATCH")
	// Delete a specific book
	v1.HandleFunc("/books/{id}", app.deleteBookHandler).Methods("DELETE")
	return v1
}
