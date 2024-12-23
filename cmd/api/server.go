package main

import (
	"fmt"
	"net/http"
	"time"
)

func (app *application) serve() error {
	// Declare a http server with some sensible timeout settings, which listens on the
	// port provided in the  config struct and uses the servemux we created above as the handler
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  app.config.env,
	})

	return srv.ListenAndServe()
}
