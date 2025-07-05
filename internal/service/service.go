package service

type Repository interface{}

type Service struct {
}

func NewService(repo Repository) *Service {
	return &Service{}
}
