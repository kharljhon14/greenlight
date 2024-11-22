package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Application version number
const version = "1.0.0"

// Configuration settings
type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *log.Logger
}

func main() {
	// Declare config instance
	var cfg config

	// Read the value of the port and env command-line flag into the config struct.
	// We default to using port number 4000 and enviroment "development"
	// If no corresponding flags are provided
	flag.IntVar(&cfg.port, "port", 4000, "API Server port")
	flag.StringVar(&cfg.env, "env", "development", "Enviroment (development|staging|production)")
	flag.Parse()

	// Initialize a new logger which writes messages to the standard out stream
	// Prefixed with the current date and time
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Declare an instance of the application struct, containing the config and the logger
	app := &application{
		config: cfg,
		logger: logger,
	}

	// Declare new servermux
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)

	// Declare a http server with some sensible timeout settings, which listens on the
	// port provided in the  config struct and uses the servemux we created above as the handler
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start http server
	logger.Printf("Starting %s server on %s", cfg.env, srv.Addr)
	err := srv.ListenAndServe()
	logger.Fatal(err)

}
