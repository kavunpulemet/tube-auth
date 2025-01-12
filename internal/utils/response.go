package utils

import (
	"encoding/json"
	"net/http"
)

const (
	ContentType     = "Content-Type"
	ApplicationJSON = "application/json"
)

type StatusResponse struct {
	Status string `json:"status"`
}

func WriteResponse(w http.ResponseWriter, statusCode int, data any) error {
	w.Header().Set(ContentType, ApplicationJSON)
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(data)
}
