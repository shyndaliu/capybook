package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/shyndaliu/capybook/pkg/capybook/model"
	"github.com/shyndaliu/capybook/pkg/capybook/validator"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "Hello! Welcome to Capybook API",
	}
	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string   `json:"title"`
		Author      string   `json:"author"`
		Year        int32    `json:"year"`
		Description string   `json:"description"`
		Genres      []string `json:"genres"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	book := &model.Book{
		Title:       input.Title,
		Author:      input.Author,
		Year:        input.Year,
		Description: input.Description,
		Genres:      input.Genres,
	}

	v := validator.New()
	if model.ValidateBook(v, book); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Books.Insert(book)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/api/v1/books/%d", book.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"book": book}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) listBooksHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string
		Author string
		Genres []string
		model.Filters
	}
	v := validator.New()
	qs := r.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Author = app.readString(qs, "author", "")
	input.Genres = app.readCSV(qs, "genres", []string{})
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.Limit = app.readInt(qs, "limit", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "title", "year", "author", "-id", "-title", "-year", "author"}

	model.ValidateFilters(v, input.Filters)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Call the GetAll() method to retrieve the movies, passing in the various filter
	// parameters.
	books, err := app.models.Books.GetAll(input.Title, input.Author, input.Genres, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Send a JSON response containing the movie data.
	err = app.writeJSON(w, http.StatusOK, envelope{"books": books}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) getBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	book, err := app.models.Books.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) updateBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	book, err := app.models.Books.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	var input struct {
		Title       *string  `json:"title"`
		Author      *string  `json:"author"`
		Year        *int32   `json:"year"`
		Description *string  `json:"description"`
		Genres      []string `json:"genres"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if input.Title != nil {
		book.Title = *input.Title
	}
	if input.Author != nil {
		book.Author = *input.Author
	}
	if input.Year != nil {
		book.Year = *input.Year
	}
	if input.Description != nil {
		book.Description = *input.Description
	}
	if input.Genres != nil {
		book.Genres = input.Genres
	}

	newbook, err := app.models.Books.Update(book)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"movie": newbook}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	err = app.models.Books.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "movie successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
