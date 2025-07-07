package service

import (
	"context"

	"micro-blog/internal/model"
)

type PostRepository interface {
	CreatePost(post *model.Post) (*model.Post, error)
	GetListPost() ([]*model.Post, error)
	LikePost(like *model.Like) error
}

type PostService struct {
	postRepo PostRepository
	userRepo UserRepository
}

func NewPostService(pr PostRepository, up UserRepository) *PostService {
	return &PostService{
		postRepo: pr,
		userRepo: up,
	}
}

func (s *PostService) CreatePost(ctx context.Context, post *model.Post) (*model.Post, error) {
	if _, err := s.userRepo.GetUserById(post.AuthorID); err != nil {
		return nil, err
	}

	return s.postRepo.CreatePost(post)
}

func (s *PostService) GetListPost(ctx context.Context) ([]*model.Post, error) {
	return s.postRepo.GetListPost()
}

func (s *PostService) LikePost(ctx context.Context, like *model.Like) error {
	if _, err := s.userRepo.GetUserById(like.UserID); err != nil {
		return err
	}

	return s.postRepo.LikePost(like)
}
