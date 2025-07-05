package response

import (
	"encoding/json"
	"net/http"

	"micro-blog/internal/handler/dto"
)

func WriteError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(dto.ErrorResponse{Message: message})
}
