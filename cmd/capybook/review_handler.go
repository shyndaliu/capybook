package main

import (
	"errors"
	"net/http"

	"github.com/shyndaliu/capybook/pkg/capybook/model"
	"github.com/shyndaliu/capybook/pkg/capybook/validator"
)

func (app *application) postReviewHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Content string `json:"content"`
		Rating  int    `json:"rating"`
	}
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	review := &model.Review{
		Content:  input.Content,
		Rating:   input.Rating,
		BookId:   id,
		AuthorId: 1,
		//rewrite
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
	review.BookTitle = book.Title

	// user, err = app.models.Users.GetByID(1)
	// if err != nil {
	// 	switch {
	// 	case errors.Is(err, model.ErrRecordNotFound):
	// 		app.notFoundResponse(w, r)
	// 	default:
	// 		app.serverErrorResponse(w, r, err)
	// 	}
	// 	return
	// }
	//review.AuthorUsername = user.Username
	review.AuthorUsername = "uldana"

	v := validator.New()
	if model.ValidateReview(v, review); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Reviews.Insert(review)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"review": review}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) listReviewsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	_, err = app.models.Books.Get(id)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		model.Filters
	}
	v := validator.New()
	qs := r.URL.Query()

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.Limit = app.readInt(qs, "limit", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "-created_at")
	input.Filters.SortSafelist = []string{"created_at", "-created_at", "rating", "-rating"}

	model.ValidateFilters(v, input.Filters)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	reviews, err := app.models.Reviews.GetAll(id, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Send a JSON response containing the movie data.
	err = app.writeJSON(w, http.StatusOK, envelope{"reviews": reviews}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) updateReviewHandler(w http.ResponseWriter, r *http.Request) {
	book_id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	review, err := app.models.Reviews.Get(book_id, 1)
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
		Content *string `json:"content"`
		Rating  *int    `json:"rating"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	review.Content = *input.Content
	review.Rating = *input.Rating

	v := validator.New()
	model.ValidateReview(v, review)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Reviews.Update(review)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"review": review}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
func (app *application) deleteReviewHandler(w http.ResponseWriter, r *http.Request) {
	book_id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	err = app.models.Reviews.Delete(book_id, 1)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "review successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
