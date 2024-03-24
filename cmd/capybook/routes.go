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
	// List books
	v1.HandleFunc("/books", app.listBooksHandler).Methods("GET")
	//Get specific book
	v1.HandleFunc("/books/{id}", app.getBookHandler).Methods("GET")
	// Update a specific book
	v1.HandleFunc("/books/{id}", app.updateBookHandler).Methods("PATCH")
	// Delete a specific book
	v1.HandleFunc("/books/{id}", app.deleteBookHandler).Methods("DELETE")

	//Users
	//Register new user
	v1.HandleFunc("/users", app.registerUserHandler).Methods("POST")
	//Get specific user
	v1.HandleFunc("/users/{username}", app.getUserHandler).Methods("GET")
	//Change the password
	v1.HandleFunc("/users/{username}", app.updateUserHandler).Methods("PATCH")
	//Delete the user
	v1.HandleFunc("/users/{username}", app.deleteUserHandler).Methods("DELETE")

	//Activate new user
	v1.HandleFunc("/users/activated", app.activateUserHandler).Methods("PUT")
	//Authenticate new user
	v1.HandleFunc("/token/authentication", app.createAuthTokenHandler).Methods("POST")

	//Reviews
	//Post new review
	v1.HandleFunc("/books/{id}/reviews", app.postReviewHandler).Methods("POST")
	//List all reviews under the book
	v1.HandleFunc("/books/{id}/reviews", app.listReviewsHandler).Methods("GET")
	//Update review details
	v1.HandleFunc("/books/{id}/reviews", app.updateReviewHandler).Methods("PATCH")
	//Delete review
	v1.HandleFunc("/books/{id}/reviews", app.deleteReviewHandler).Methods("DELETE")
	return v1
}
