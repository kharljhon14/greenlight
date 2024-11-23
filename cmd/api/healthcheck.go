package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// Map
	data := map[string]string{
		"status":     "available",
		"enviroment": app.config.env,
		"version":    version,
	}

	// Pass  the json map to the json.marshal
	js, err := json.Marshal(data)
	if err != nil {
		// Send generic message
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		return
	}

	// Append newline to the json
	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")

	w.Write(js)
}
