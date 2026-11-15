package main

import (
	"errors"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Call run() to start the application.
// If run() returns an error, we log the error and exit the application with status code 1 (indicating failure).
func main() {
	if err := run(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

// func run() contains the main application setup.
// Separated logic from main() so that the application is easier to organize, test, and maintain.
// It returns an error if something goes wrong during startup or while running the HTTP server.
func run() error {

	// Load environment variables from the .env file into the application's environment.
	// godotenv.Load() reads the .env file and makes its key-value pairs available as environment variables that can be accessed using os.Getenv().
	// This allows us to keep configuration values, such as database credentials, outside the source code instead of hardcoding them.
	if err := godotenv.Load(); err != nil {
		return err
	}

	// Read the DATABASE_URL environment variable.
	// os.Getenv() retrieves the value of DATABASE_URL from the application's environment.
	// This value was made available by godotenv.Load().
	// Example .env entry:
	// DATABASE_URL="postgres://user:password@localhost:5432/mydb"
	url := os.Getenv("DATABASE_URL")

	// Make sure the DATABASE_URL environment variable was provided.
	// If it is empty, we cannot establish a database connection, so we return an error and stop application startup.
	if url == "" {
		return errors.New("DATABASE_URL is required")
	}

	// Attempt to establish a connection to the database.
	// Receive the database connection URL and return:
	// 1. db: The database connection.
	// 2. err: An error if the connection could not be established.
	db, err := openDB(url)

	// Check whether the database connection failed.
	// If it did, return the error so main() can handle it.
	if err != nil {
		return err
	}

	// defer schedules db.Close() to run when run() finishes.
	// This ensures we release the database connection when the application shuts down to avoid hogging resources.
	defer db.Close()

	// HTTP Methods/Actions (CRUD Requests), HTTP Paths (Requested address residing within the server), HTTP Handler (Function that processes the request), HTTP Status (Response code), HTTP Router (Matches: Request -> Path -> Handler)
	// Create a new HTTP request multiplexer (mux).
	// A mux acts as the application's router. It examines incoming HTTP requests and directs them to the appropriate handler based on the HTTP method and URL path.
	// For example:
	// GET  /health       -> Health check handler
	// POST /api/preview  -> Preview handler

	mux := http.NewServeMux()

	// Register a handler for GET requests to /health.
	// mux.HandleFunc() associates a URL pattern with a handler function.
	// When an incoming request matches the registered pattern, Go calls the associated handler (healthHandler in this case).
	// The handler function (healthHandler) called, is being implemented separately in health.go.
	mux.HandleFunc("GET /health", healthHandler(db))

	// Register a handler for POST requests to /api/preview.
	// When a client sends a POST request to this endpoint, Go calls previewHandler() to process the request.

	// The handler is defined separately in handlers.go.
	// Keeping handlers in separate files helps organize the application as it grows in size & complexity.
	mux.HandleFunc("POST /api/preview", previewHandler)

	// Create an HTTP server configuration.
	// http.Server allows us to configure how our HTTP server behaves.
	server := &http.Server{
		// Addr specifies the network address and port where the HTTP server will listen for requests.
		Addr: "localhost:8080",

		// Handler tells the server which router should handle incoming HTTP requests.
		// Here, we pass our mux so it can direct requests to the correct handler.
		Handler: mux,

		// ReadHeaderTimeout limits how long the server will wait to receive the request headers from a client.
		// This helps protect the server against clients that connect but take too long to send their headers.
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Print a message to the console so we know the API is about to start listening on localhost:8080.
	log.Println("API running at http://localhost:8080")

	// Start the HTTP server.
	// ListenAndServe() begins listening for incoming HTTP requests on the address configured above.

	// It uses our mux to route incoming requests to their handlers.

	// This function blocks while the server is running.
	// If the server encounters an error, it returns that error.
	return server.ListenAndServe()
}

