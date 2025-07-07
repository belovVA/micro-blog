package converter

import (
	"github.com/google/uuid"
	"micro-blog/internal/handler/dto"
	"micro-blog/internal/model"
)

func ToPostModelFromReq(req *dto.CreatePostReq) (*model.Post, error) {
	authorID, err := uuid.Parse(req.AuthorID)
	if err != nil {
		return nil, err
	}

	return &model.Post{
		ID:       uuid.Nil,
		AuthorID: authorID,
		Text:     req.Text,
	}, nil
}

func ToPostRespFromModel(post *model.Post) *dto.PostResp {
	return &dto.PostResp{
		ID:       post.ID.String(),
		AuthorID: post.AuthorID.String(),
		Text:     post.Text,
		Likes:    post.Likes,
	}
}
