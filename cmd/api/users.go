package main

import (
	"errors"
	"net/http"

	"github.com/thecodephilic-guy/greenlight/internal/data"
	"github.com/thecodephilic-guy/greenlight/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	//Anonymous struct to hold the expected data from the request body:
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Copy the data from the request body into a new User struct. Notice also that we
	// set the Activated field to false, which isn't strictly necessary because the
	// Activated field will have the zero-value of false by default. But setting this
	// explicitly helps to make our intentions clear to anyone reading the code.
	user := &data.Users{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	//Use the Password.set() method to generate and store the hashed and plaintext passwords
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	//validate the user struct and return the error messages to the client if any of
	//the checks fail.
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	//Insert the user data into db
	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		// If we get a ErrDuplicateEmail error, use the v.AddError() method to manually
		// add a message to the validator instance, and then call our
		// failedValidationResponse() helper.
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Launch a goroutine which runs an anonymous function that sends the welcome email.
	go func() {
		err = app.mailer.Send(user.Email, "user_welcome.html", user)
		if err != nil {
			// Importantly, if there is an error sending the email then we use the
			// app.logger.PrintError() helper to manage it, instead of the
			// app.serverErrorResponse()
			app.logger.PrintError(err, nil)
		}
	}()

	// Note that we also change this to send the client a 202 Accepted status code.
	// This status code indicates that the request has been accepted for processing, but
	// the processing has not been completed.
	err = app.writeJSON(w, http.StatusAccepted, envelop{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
