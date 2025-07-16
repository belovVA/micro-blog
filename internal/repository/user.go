package repository

import (
	"sync"

	"github.com/google/uuid"
	"micro-blog/internal/model"
)

type UserRepo struct {
	Users map[string]*model.User
	mu    sync.RWMutex
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		Users: make(map[string]*model.User),
		mu:    sync.RWMutex{},
	}
}

func (r *UserRepo) CreateUser(user *model.User) (*model.User, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	user.ID = id

	r.mu.Lock()
	defer r.mu.Unlock()
	r.Users[user.Name] = user

	return user, nil
}

func (r *UserRepo) GetUserByName(name string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if val, ok := r.Users[name]; ok {
		return val, nil
	}

	return nil, model.ErrUserNotFound
}

func (r *UserRepo) GetUserById(id uuid.UUID) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, val := range r.Users {
		if val.ID == id {
			return val, nil
		}
	}

	return nil, model.ErrUserNotFound
}
