package repository

import (
	"github.com/google/uuid"
	"micro-blog/internal/model"
)

type UserRepo struct {
	Users map[string]*model.User
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		Users: make(map[string]*model.User),
	}
}

func (r *UserRepo) CreateUser(user *model.User) (*model.User, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	user.ID = id
	r.Users[user.Name] = user

	return user, nil
}

func (r *UserRepo) GetUserByName(name string) (*model.User, error) {
	if val, ok := r.Users[name]; ok {
		return val, nil
	}
	return nil, model.ErrUserNotFound
}
