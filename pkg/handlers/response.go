package handlers

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	errors []error
}

func writeResponse(w http.ResponseWriter, statusCode int, errors []string) {
	w.WriteHeader(statusCode)

	if len(errors) > 0 {
		jw := json.NewEncoder(w)
		jw.Encode(errors)
	}
}
