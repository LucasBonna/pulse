package utils

import (
	"encoding/json"
	"net/http"
)

type GenericResponse struct {
	Data any `json:"data"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

func WriteJsonResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	genericResp := GenericResponse{
		Data: data,
	}

	json.NewEncoder(w).Encode(genericResp)
}

func WriteJsonError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResp := ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	}

	json.NewEncoder(w).Encode(errorResp)
}
