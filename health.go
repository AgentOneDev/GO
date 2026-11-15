package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// healthHandler checks whether the application and its database are running.
// It accepts a database connection and returns an HTTP handler function.
//
// Parameters:
// db *sql.DB: A database connection pool used to communicate with the database.
//
// Returns:
// http.HandlerFunc: A function that handles incoming HTTP requests and writes HTTP responses.
func healthHandler(db *sql.DB) http.HandlerFunc {

	// Return an anonymous handler function that matches the http.HandlerFunc type.
	// This allows the handler to access the database connection (db) through its closure.
	return func(w http.ResponseWriter, r *http.Request) {

		// Create a context with a 5-second timeout.
		// r.Context() provides the context associated with the incoming HTTP request.
		// context.WithTimeout() adds a time limit to operations performed using this context.
		// If the database query takes longer than 5 seconds, the context is canceled.
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)

		// Ensure the context is canceled when the handler finishes.
		// This releases resources associated with the context.
		defer cancel()

		// Declare an integer variable to store the result of the database query.
		var result int

		// Execute a database query using the context with a 5-second timeout.
		// QueryRowContext() executes a query expected to return a single row.
		// Scan() copies the returned value into the result variable.
		// SELECT 1 is a simple query used to verify that the database is reachable
		// and can successfully execute a query.
		err := db.QueryRowContext(ctx, "SELECT 1").Scan(&result)

		// Check whether the database query returned an error.
		if err != nil {

			// Log the database health check error.
			log.Printf("health query: %v", err)

			// Send an HTTP 503 Service Unavailable response to the client.
			http.Error(w, "database unavailable", http.StatusServiceUnavailable)

			// Stop executing the handler after sending the error response.
			return
		}

		// Set the Content-Type response header to indicate that the response
		// body will contain JSON data.
		w.Header().Set("Content-Type", "application/json")

		// Create a JSON encoder and use it to write the response directly to the client.
		// json.NewEncoder(w) creates an encoder that writes JSON to the HTTP response writer.
		// Encode() converts the Go map into JSON and writes it to the response.
		if err := json.NewEncoder(w).Encode(map[string]string{
			"status":   "ok",
			"database": "connected",
		}); err != nil {

			// Log an error if the JSON response cannot be encoded or written.
			log.Printf("encode health: %v", err)
		}
	}
}

