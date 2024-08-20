package main

import "net/http"

func (app *application) routes() http.Handler {
	router := http.NewServeMux()

	base := CreateStack(app.recoverPanic, app.rateLimit, app.authenticate)
	protected := CreateStack(app.requireActivatedUser)

	router.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)

	// Movies
	router.HandleFunc("GET /v1/movies", protected.ToHandlerFunc(app.listMoviesHandler))
	router.HandleFunc("GET /v1/movies/{id}", protected.ToHandlerFunc(app.showMovieHandler))
	router.HandleFunc("POST /v1/movies", protected.ToHandlerFunc(app.createMovieHandler))
	router.HandleFunc("PATCH /v1/movies/{id}", protected.ToHandlerFunc(app.updateMovieHandler))
	router.HandleFunc("DELETE /v1/movies/{id}", protected.ToHandlerFunc(app.deleteMovieHandler))

	// Users
	router.HandleFunc("POST /v1/users", app.registerUserHandler)
	router.HandleFunc("PUT /v1/users/activated", app.activateUserHandler)

	// Tokens
	router.HandleFunc("POST /v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return base(router)
}
