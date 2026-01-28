package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/thecodephilic-guy/greenlight/internal/data"
	"github.com/thecodephilic-guy/greenlight/internal/validator"
)

// "POST /v1/movies"
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	// Declare an anonymous struct to hold the information that we expect to be in the
	// HTTP request body (note that the field names and types in the struct are a subset
	// of the Movie struct that we created earlier). This struct will be our *target
	// decode destination*.
	var input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Copy the values from the input struct to a new Movie struct.
	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	// Initialize a new Validator instance:
	v := validator.New()

	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	//Now calling the Insert() method on our movie model, passing in a pointer to the
	//validated movie struct. This will create a record in the database and update the
	//movie struct with the system-generated information.
	err = app.model.Movies.Insert(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	// When sending a HTTP response, we want to include a Location header to let the
	// client know which URL they can find the newly-created resource at. We make an
	// empty http.Header map and then use the Set() method to add a new Location header,
	// interpolating the system-generated ID for our new movie in the URL.
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	// Write a JSON response with a 201 Created status code, the movie data in the
	// response body, and the Location header.
	err = app.writeJSON(w, http.StatusCreated, envelop{"movie": movie}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// "GET /v1/movie/:id"
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIdParams(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.model.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	//Encoding the struct into JSON using helper funtion and sending reponse
	err = app.writeJSON(w, http.StatusOK, envelop{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// "PUT /v1/movies/:id"
func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIdParams(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	//Fetching the record from DB to make sure it exists:
	movie, err := app.model.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	//input struct to hold the expected data:
	var input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	//copy values to movie record:
	movie.Title = input.Title
	movie.Year = input.Year
	movie.Runtime = input.Runtime
	movie.Genres = input.Genres

	// Validate the updated movie record, sending the client 1 422 Unprocessable Entity
	// response if any checks fail.
	v := validator.New()

	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Now update the DB:
	err = app.model.Movies.Update(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	//send response to the client:
	err = app.writeJSON(w, http.StatusOK, envelop{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// "DELETE /v1/movies/:id"
func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIdParams(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	//Deleting the movie data from the databse and sending 404 not found if not present:
	err = app.model.Movies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)

		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	//sending the response to the client:
	err = app.writeJSON(w, http.StatusOK, envelop{"message": "movie successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
