package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) createMovieHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Create a new movie...\n")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Show the details of movie %d\n", id)
}
