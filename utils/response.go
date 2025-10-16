package utils

import (
    "encoding/json"
    "net/http"
)

type APIResponse struct {
    Status  bool        `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

func Success(w http.ResponseWriter, message string, data interface{}, code int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    _ = json.NewEncoder(w).Encode(APIResponse{Status: true, Message: message, Data: data})
}

func Error(w http.ResponseWriter, message string, code int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    _ = json.NewEncoder(w).Encode(APIResponse{Status: false, Message: message})
}


