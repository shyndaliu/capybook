package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/shyndaliu/capybook/pkg/capybook/model"
)

func (app *application) authenticate(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Vary", "Authorization")
		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			r = app.contextSetUser(r, model.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}
		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}
		token := headerParts[1]

		if r.URL.Path == "/api/v1/token/refresh" {
			fmt.Print(token)
			userRefresh, err := app.auth.ValidateRefreshToken(token)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			user, err := app.models.Users.GetByUsername(userRefresh.Username)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			actualCustomKey := app.auth.GenerateCustomKey(user.Username, user.TokenHash)
			if userRefresh.CustomKey != actualCustomKey {
				app.invalidAuthenticationTokenResponse(w, r)
				return
			}
			r = app.contextSetUser(r, user)

		} else {
			userAccess, err := app.auth.ValidateAccessToken(token)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}
			user, err := app.models.Users.GetByUsername(userAccess.Username)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}
			r = app.contextSetUser(r, user)

		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)
		if user.IsAnonymous() {
			app.authenticationRequiredResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Checks that a user is both authenticated and activated.
func (app *application) requireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
	// Rather than returning this http.HandlerFunc we assign it to the variable fn.
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)
		// Check that a user is activated.
		if !user.Activated {
			app.inactiveAccountResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
	// Wrap fn with the requireAuthenticatedUser() middleware before returning it.
	return app.requireAuthenticatedUser(fn)
}

func (app *application) requirePermission(code string, next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)
		permissions, err := app.models.Permissions.GetAllForUser(user.ID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		if !permissions.Include(code) {
			app.notPermittedResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	}
	return app.requireActivatedUser(fn)
}
