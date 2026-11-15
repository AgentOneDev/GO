package main

import (
	"io"
	"encoding/json"
	"net/http"
	"log"
	"mime"
)

func readAuthJSON(w http.ResponseWriter, r *http.Request, destination any) bool {
	contentType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || contentType != "application/json" {
		authError(w, http.StatusUnsupportedMediaType, "send application/json")
		return false
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)

	if err != nil {
		authError(w, http.StatusBadRequest, "could not read request body")
		return false
	}

	err = json.Unmarshal(body, destination)

	if err != nil {
		authError(w, http.StatusBadRequest, "invalid json")
		return false
	}
	return true
}

func authJson(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().set("Cache-Control", "no-store")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(value)

	if err != nil {
		log.Printf("write auth response: %v", err)
	}
}
func authError(w http:ResponseWriter, status int, message string) {
	authJson(w, status, map[string]string{"error", message})
}
