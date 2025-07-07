package model

import "github.com/google/uuid"

type Post struct {
	ID       uuid.UUID
	AuthorID uuid.UUID
	Text     string
	Likes    []uuid.UUID
}
