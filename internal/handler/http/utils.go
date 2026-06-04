package http

import (
	"encoding/json"
	"net/http"
)

// Errors
type ErrorResponse struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, ErrorResponse{
		Error:  msg,
		Status: status,
	})
}
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
