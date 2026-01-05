package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/thecodephilic-guy/greenlight/internal/data"
)

// "POST /v1/movies"
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new movie")
}

// "GET /v1/movie/:id"
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIdParams(r)

	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Create a new instance of the Movie struct, containing the ID we extracted from
	// the URL and some dummy data. Also notice that we deliberately haven't set a
	// value for the Year field.”

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	//Encoding the struct into JSON using helper funtion and sending reponse
	err = app.writeJSON(w, http.StatusOK, envelop{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
