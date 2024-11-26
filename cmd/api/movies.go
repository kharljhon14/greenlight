package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kharljhon14/greenlight/internal/data"
)

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	//  Interpolated URL parameters will be
	// stored in the request context.
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundErrorResponse(w, r)
		return
	}

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Hello",
		Runtime:   132,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
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
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
	}

	// Dump the contens of the input struct in a http response
	fmt.Fprintf(w, "%+v\n", input)
}
