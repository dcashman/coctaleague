package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Define a config struct to hold all the configuration settings for our application.
type config struct {
	port int    // Network port on which to listen.
	env  string // Name of current operating environment.
}

// Define an application struct to hold the dependencies for our HTTP handlers,
// helpers, and middleware. Also useful for testing.
type application struct {
	config   config
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// Declare an instance of the config struct.
	var cfg config

	// Read the value of the port and env command-line flags into the config struct. We
	// default to using the port number 4000 and the environment "development" if no
	// corresponding flags are provided.
	flag.IntVar(&cfg.port, "port", 4000, "Server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Our applicaiton struct: used for passing dependencies around neatly.
	app := &application{
		config:   cfg,
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// Declare a HTTP server with some sensible timeout settings, which listens on the
	// port provided in the config struct and uses the servemux we created above as the
	// handler.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start the HTTP server.
	infoLog.Printf("Starting server on %d", app.config.port)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
