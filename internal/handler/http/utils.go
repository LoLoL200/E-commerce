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
	w.Header().Set("content-Type", "app;ication/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
