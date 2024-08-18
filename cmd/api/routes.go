package main

import "net/http"

func (app *application) routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)

	// Movies
	router.HandleFunc("GET /v1/movies", app.listMoviesHandler)
	router.HandleFunc("GET /v1/movies/{id}", app.showMovieHandler)
	router.HandleFunc("POST /v1/movies", app.createMovieHandler)
	router.HandleFunc("PATCH /v1/movies/{id}", app.updateMovieHandler)
	router.HandleFunc("DELETE /v1/movies/{id}", app.deleteMovieHandler)

	// Users
	router.HandleFunc("POST /v1/users", app.registerUserHandler)
	router.HandleFunc("PUT /v1/users/activated", app.activateUserHandler)

	// Tokens
	router.HandleFunc("POST /v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
