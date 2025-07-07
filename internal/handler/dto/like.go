package dto

type LikeRequest struct {
	UserID string `json:"user_id" validate:"required"`
}
