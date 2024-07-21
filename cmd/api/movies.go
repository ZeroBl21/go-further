package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ZeroBl21/go-further/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Create a new movie...\n")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	movie := data.Movie{
		ID:      id,
		Title:   "Casablanca",
		Runtime: 102,
		Genres:  []string{"drama", "romance", "war"},

		CreatedAt: time.Now(),
		Version:   1,
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil); err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and coult not process your request",
			http.StatusInternalServerError)
	}
}
