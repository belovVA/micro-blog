package service

import (
	"context"

	"micro-blog/internal/model"
)

type UserRepository interface {
	CreateUser(user *model.User) (*model.User, error)
	GetUserByName(name string) (*model.User, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Authenticate(ctx context.Context, user *model.User) (*model.User, error) {
	if u, err := s.repo.GetUserByName(user.Name); err == nil {
		return u, nil
	}

	var err error
	if user, err = s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}
