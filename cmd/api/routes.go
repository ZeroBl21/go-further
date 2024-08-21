package main

import "net/http"

func (app *application) routes() http.Handler {
	router := http.NewServeMux()

	base := CreateStack(app.recoverPanic, app.rateLimit, app.authenticate)

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

	// Tokens
	router.HandleFunc("POST /v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return base(router)
}
