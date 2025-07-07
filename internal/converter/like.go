package converter

import (
	"github.com/google/uuid"
	"micro-blog/internal/handler/dto"
	"micro-blog/internal/model"
)

func ToLikeModelFromReq(req *dto.LikeRequest, postIDStr string) (*model.Like, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, err
	}

	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		return nil, err
	}

	return &model.Like{
		UserID: userID,
		PostID: postID,
	}, nil
}
