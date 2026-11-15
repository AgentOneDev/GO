package main

import (
	"context"
	"database/sql" // Import the database/sql package for database operations
	"fmt"
	"time"
)

// openDb opens a database connection using the provided URL and returns a pointer to the sql.DB object and an error if any occurs during the connection process.
func openDB(url string) (*sql.DB, error) {
	db, err := sql.Open("pgx", url) // Open a connection to the database using the pgx driver and the provided URL
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second) // Create a context with a timeout of 15 seconds to ensure that the database connection is established within this time frame
	defer cancel()                                                           // Ensure that the context is canceled after the function returns to free up resources

	if db.PingContext(ctx); err != nil {
		db.Close()                                                 // Close the database connection if the ping fails
		return nil, fmt.Errorf("failed to ping database: %w", err) // Return an error if the ping fails
	}
	return db, nil // Return db pool if connection is successful
}

