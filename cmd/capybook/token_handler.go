package main

import (
	"errors"
	"net/http"

	"github.com/shyndaliu/capybook/pkg/capybook/model"
	"github.com/shyndaliu/capybook/pkg/capybook/validator"
)

func (app *application) createAuthTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	v := validator.New()
	model.ValidatePasswordPlaintext(v, input.Password)
	model.ValidateEmailOrUsername(v, input.Username, input.Email)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	model.ValidateUsername(v, input.Username)
	usernameValid := v.Valid()
	model.ValidateEmail(v, input.Email)
	emailValid := v.Valid()
	if !usernameValid && !emailValid {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	var user *model.User
	if emailValid {
		user, err = app.models.Users.GetByEmail(input.Email)
	} else {
		user, err = app.models.Users.GetByUsername(input.Username)
	}
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}
	accessToken, err := app.auth.GenerateAccessToken(user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	refreshToken, err := app.auth.GenerateRefreshToken(user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusCreated, envelope{"access_token": accessToken, "refresh_token": refreshToken}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) refreshTokenandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}
	accessToken, err := app.auth.GenerateAccessToken(user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusCreated, envelope{"access_token": accessToken}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
