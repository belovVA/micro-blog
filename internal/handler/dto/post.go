package dto

import "github.com/google/uuid"

type CreatePostReq struct {
	AuthorID string `json:"author_id" validate:"required"`
	Text     string `json:"text"`
}

type PostResp struct {
	ID       string      `json:"id"`
	AuthorID string      `json:"author_id"`
	Text     string      `json:"text"`
	Likes    []uuid.UUID `json:"likes"`
}
