package dto

type CreateUserReq struct {
	Name string `json:"name" validate:"required"`
}

type CreateUserResp struct {
	ID string `json:"id"`
}
