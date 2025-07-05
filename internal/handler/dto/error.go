package dto

type ErrorResponse struct {
	// Сообщение об ошибке, описывающее проблему.
	Message string `json:"message,omitempty"`
}
