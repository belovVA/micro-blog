package converter

import (
	"github.com/google/uuid"
	"micro-blog/internal/handler/dto"
	"micro-blog/internal/model"
)

func ToUserModelFromReq(req *dto.CreateUserReq) *model.User {
	return &model.User{
		ID:   uuid.Nil,
		Name: req.Name,
	}
}

func ToUserRespFromModel(user *model.User) *dto.CreateUserResp {
	return &dto.CreateUserResp{
		ID: user.ID.String(),
	}
}
