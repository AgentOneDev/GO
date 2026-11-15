package main

import (
	"encoding/json" 
	"log"           
	"net/http"      
	"strings"       
)

// previewInput defines the expected structure of the JSON input received by the previewHandler function.
//
// Fields:
// Title: The title provided by the client.
// Body: The body content provided by the client.
//
// Struct tags (`json:"title"` and `json:"body"`) specify how JSON properties map to the corresponding Go struct fields.
type previewInput struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// previewHandler handles incoming HTTP requests containing JSON data.
// It validates the request and returns the decoded data as a JSON response.
//
// Parameters:
// w http.ResponseWriter: Used to write the HTTP response back to the client.
// r *http.Request: Contains information about the incoming HTTP request.
func previewHandler(w http.ResponseWriter, r *http.Request) {

	// Declare a variable of type previewInput to store the decoded JSON data.
	var input previewInput

	// Decode the incoming JSON request body into the input variable.
	// json.NewDecoder(r.Body) creates a decoder that reads JSON from the request body.
	// Decode(&input) converts the JSON data into the corresponding Go struct fields.
	// The & operator passes a pointer so Decode() can modify the input variable.
	err := json.NewDecoder(r.Body).Decode(&input)

	// Check whether the JSON decoding returned an error.
	// If the request body contains invalid JSON, return a 400 Bad Request response.
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Remove leading and trailing whitespace from the Title field.
	// This prevents titles containing only spaces from being accepted especially for user input validation.
	input.Title = strings.TrimSpace(input.Title)

	// Validate that the Title field is not empty.
	// If the title is empty after trimming whitespace, return a 400 Bad Request response.
	if input.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	// Set the Content-Type response header to indicate that the response body will contain JSON data.
	w.Header().Set("Content-Type", "application/json")

	// Encode the input struct back into JSON and write it directly to the HTTP response.
	// json.NewEncoder(w) creates an encoder that writes JSON to the response writer.
	// Encode(input) converts the Go struct into JSON.
	if err := json.NewEncoder(w).Encode(input); err != nil {

		// Log an error if the JSON response cannot be encoded or written.
		log.Printf("Error encoding JSON: %v", err)
	}
}

