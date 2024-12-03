package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Application version number
const version = "1.0.0"

// Configuration settings
type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
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

	// Read the DSN value from the db-dsn command line flag into the config struct
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://greenlight:password@localhost/greenlight?sslmode=disable", "PostgreSQL DSN")

	flag.Parse()

	// Initialize a new logger which writes messages to the standard out stream
	// Prefixed with the current date and time
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	// Defer a call to db.Close() so that the connection pool is closed before the main() function exits
	defer db.Close()

	logger.Printf("database connection pool established")

	// Declare an instance of the application struct, containing the config and the logger
	app := &application{
		config: cfg,
		logger: logger,
	}

	// Declare a http server with some sensible timeout settings, which listens on the
	// port provided in the  config struct and uses the servemux we created above as the handler
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start http server
	logger.Printf("Starting %s server on %s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)

}

func openDB(cfg config) (*sql.DB, error) {
	// Create a empty connection pool using the DSN from the config
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// Create contect with 5 second timeout deadline
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Establish a new connection to the database
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
