package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/kharljhon14/greenlight/internal/data"
	"github.com/kharljhon14/greenlight/internal/validator"
)

func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {
	// Placeholder for the expected input values
	var input struct {
		Title    string
		Genres   []string
		Page     int
		PageSize int
		Sort     string
	}

	// New Validator
	v := validator.New()

	qs := r.URL.Query()

	// Extract the title and genres
	input.Title = app.readString(qs, "title", "")
	input.Genres = app.readCSV(qs, "genres", []string{})

	// Extract page and pagesize
	input.Page = app.readInt(qs, "page", 1, v)
	input.PageSize = app.readInt(qs, "page_size", 20, v)

	// Extract sort
	input.Sort = app.readString(qs, "sort", "id")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	//  Interpolated URL parameters will be
	// stored in the request context.
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundErrorResponse(w, r)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundErrorResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Encode the struct to JSON and send as the HTTP response
	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		// app.logger.Println(err)
		// http.Error(w, "Ther server encountered a problem and could not process your request", http.StatusInternalServerError)
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {

	// Struct for holding information we expect to be in the http request
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
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

	// Init new validator
	v := validator.New()

	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Movies.Insert(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	// Set location header for the newly created movie
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	// Dump the contens of the input struct in a http response
	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the movie ID from the URL
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundErrorResponse(w, r)
		return
	}

	// Fetch the existing movie
	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundErrorResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Input struct for holding expected data
	var input struct {
		Title   *string       `json:"title"`
		Year    *int32        `json:"year"`
		Runtime *data.Runtime `json:"runtime"`
		Genres  []string      `json:"genres"`
	}

	// Read JSON request body
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		movie.Title = *input.Title
	}

	if input.Year != nil {
		movie.Year = *input.Year
	}

	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}

	if input.Genres != nil {
		movie.Genres = input.Genres
	}

	// Validate the updated movie record
	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Pass the updated movie
	err = app.models.Movies.Update(movie)
	if err != nil {

		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Write the updated movie record
	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundErrorResponse(w, r)
		return
	}

	err = app.models.Movies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundErrorResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "movie successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
