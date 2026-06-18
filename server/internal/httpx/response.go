package httpx

import (
	"encoding/json"
	"net/http"
)

type Envelope struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func OK(w http.ResponseWriter, data interface{}) {
	Write(w, http.StatusOK, Envelope{Code: 0, Message: "ok", Data: data})
}

func Fail(w http.ResponseWriter, status int, code int, message string) {
	Write(w, status, Envelope{Code: code, Message: message, Data: nil})
}

func Write(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
