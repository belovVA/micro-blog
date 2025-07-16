package service

import (
	"context"

	"micro-blog/internal/model"
	"micro-blog/internal/queue"
)

type PostRepository interface {
	CreatePost(post *model.Post) (*model.Post, error)
	GetListPost() ([]*model.Post, error)
	LikePost(like *model.Like) error
}

type PostService struct {
	postRepo  PostRepository
	userRepo  UserRepository
	likeQueue queue.LikeEnqueuer
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

	if s.likeQueue == nil {
		return model.ErrLikeQueue
	}

	s.likeQueue.Enqueue(like)
	return nil
}

func (s *PostService) HandleLike(ctx context.Context, like *model.Like) error {
	return s.postRepo.LikePost(like)
}

func (s *PostService) AttachLikeQueue(q queue.LikeEnqueuer) {
	s.likeQueue = q
}
