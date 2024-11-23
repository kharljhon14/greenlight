package main

import (
	"fmt"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// Fixed json response
	js := `{"status":"available", "enviroment": %q, "version": %q}`
	js = fmt.Sprintf(js, app.config.env, version)

	// Set Content-Type to application/json on the header
	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte(js))
}
