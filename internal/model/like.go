package model

import "github.com/google/uuid"

type Like struct {
	UserID uuid.UUID
	PostID uuid.UUID
}
