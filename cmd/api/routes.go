package main

import (
	"expvar"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := http.NewServeMux()

	base := CreateStack(
		app.metrics,
		app.recoverPanic,
		app.enableCORS,
		app.rateLimit,
		app.authenticate,
	)

	router.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)

	// Movies
	router.HandleFunc(
		"GET /v1/movies",
		app.requirePermissions("movies:read", app.listMoviesHandler),
	)
	router.HandleFunc(
		"GET /v1/movies/{id}",
		app.requirePermissions("movies:read", app.showMovieHandler),
	)
	router.HandleFunc(
		"POST /v1/movies",
		app.requirePermissions("movies:write", app.createMovieHandler),
	)
	router.HandleFunc(
		"PATCH /v1/movies/{id}",
		app.requirePermissions("movies:write", app.updateMovieHandler),
	)
	router.HandleFunc(
		"DELETE /v1/movies/{id}",
		app.requirePermissions("movies:write", app.deleteMovieHandler),
	)

	// Users
	router.HandleFunc("POST /v1/users", app.registerUserHandler)
	router.HandleFunc("PUT /v1/users/activated", app.activateUserHandler)
	router.HandleFunc("PUT /v1/users/password", app.updateUserPasswordHandler)

	// Tokens
	router.HandleFunc("POST /v1/tokens/activation", app.createActivationTokenHandler)
	router.HandleFunc("POST /v1/tokens/authentication", app.createAuthenticationTokenHandler)
	router.HandleFunc("POST /v1/tokens/password-reset", app.createPasswordResetTokenHandler)

	// Expvar
	router.Handle("GET /debug/vars", expvar.Handler())

	return base(router)
}
