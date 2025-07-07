package response

import (
	"encoding/json"
	"net/http"
)

func SuccessJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}

func SuccessText(w http.ResponseWriter, text string, statusCode int) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(text))
}

func SuccessCode(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}
