package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponseBody struct {
	Error string `json:"error"`
}

type SuccessResponseBody struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func JSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

func ErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	response := ErrorResponseBody{Error: message}
	JSONResponse(w, response, statusCode)
}

func SuccessResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	JSONResponse(w, data, statusCode)
}