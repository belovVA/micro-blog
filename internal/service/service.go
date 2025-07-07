package service

type Repository interface {
	UserRepository
	PostRepository
}

type Service struct {
	*UserService
	*PostService
}

func NewService(repo Repository) *Service {
	return &Service{
		UserService: NewUserService(repo),
		PostService: NewPostService(repo, repo),
	}
}
