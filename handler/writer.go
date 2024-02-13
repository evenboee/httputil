package handler

import (
	"encoding/json"
	"net/http"
)

const (
	MIMEJSON = "application/json"
	MIMEText = "text/plain"
)

func WriteContentType(w http.ResponseWriter, value string) {
	w.Header().Set("Content-Type", value)
}

func JSON(w http.ResponseWriter, code int, data any) {
	WriteContentType(w, MIMEJSON)
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(data)
}

func String[T string | []byte](w http.ResponseWriter, code int, data T) {
	WriteContentType(w, MIMEText)
	w.WriteHeader(code)
	w.Write([]byte(data))
}
