package responses

import (
	"encoding/json"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

type APIResponse struct {
	StatusCode int         `json:"code"`
	Status     string      `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}

func ResponseError(w http.ResponseWriter, code int, message string, payload interface{}) {
	respondWithJSON(w, code, APIResponse{
		StatusCode: code,
		Status:     "error",
		Message:    message,
	})
}

func ResponseSuccess(w http.ResponseWriter, code int, message string, payload interface{}) {
	respondWithJSON(w, code, APIResponse{
		StatusCode: code,
		Status:     "success",
		Message:    message,
		Data:       payload,
	})
}
