package handlers

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	errors []error
}

func writeResponse(w http.ResponseWriter, statusCode int, body interface{}) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(body)
}
