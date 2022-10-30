package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

// WriteError writes the specified status code and message to the http.ResponseWriter
// using the ErrorResponse struct.
func WriteError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	errorMessage := ErrorResponse{Error: message}
	json.NewEncoder(w).Encode(errorMessage)
}
