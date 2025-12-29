package utils

import (
	"alfdwirhmn/inventory/dto"
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

func JSONSuccess(w http.ResponseWriter, code int, message string, data any) {
	resp := Response{
		Status:  true,
		Message: message,
		Data:    data,
	}
	writeJSON(w, code, resp)
}

func JSONError(w http.ResponseWriter, code int, message string, errors any) {
	resp := Response{
		Status:  false,
		Message: message,
		Errors:  errors,
	}
	writeJSON(w, code, resp)
}

func JSONWithPagination(w http.ResponseWriter, code int, message string, data any, pagination dto.Pagination) {
	resp := map[string]interface{}{
		"status":     true,
		"message":    message,
		"data":       data,
		"pagination": pagination,
	}
	writeJSON(w, code, resp)
}

// helper internal
func writeJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}
