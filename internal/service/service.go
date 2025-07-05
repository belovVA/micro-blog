package service

type Repository interface {
	UserRepository
}

type Service struct {
	*UserService
}

func NewService(repo Repository) *Service {
	return &Service{
		UserService: NewUserService(repo),
	}
}
