package http

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

type ErrorResponse struct {
	Error struct {
		Type    string `json:"type"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	RequestID string `json:"request_id,omitempty"`
	Timestamp string `json:"timestamp"`
}

func Error(w http.ResponseWriter, status int, errType, code, message string) {
	resp := ErrorResponse{}
	resp.Error.Type = errType
	resp.Error.Code = code
	resp.Error.Message = message
	JSON(w, status, resp)
}
