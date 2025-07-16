package repository

import (
	"sync"

	"github.com/google/uuid"
	"micro-blog/internal/model"
)

const (
	initPostsCapacity = 100
)

type PostRepo struct {
	Posts []*model.Post
	mu    sync.RWMutex
}

func NewPostRepo() *PostRepo {
	return &PostRepo{
		Posts: make([]*model.Post, 0, initPostsCapacity),
		mu:    sync.RWMutex{},
	}
}

func (r *PostRepo) CreatePost(post *model.Post) (*model.Post, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	post.ID = id
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Posts = append(r.Posts, post)
	return post, nil
}

func (r *PostRepo) GetListPost() ([]*model.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.Posts, nil
}

func (r *PostRepo) LikePost(like *model.Like) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, post := range r.Posts {
		if post.ID == like.PostID {
			var alreadyLiked bool
			for _, v := range post.Likes {
				if v == like.UserID {
					alreadyLiked = true
				}
			}
			if !alreadyLiked {
				r.Posts[i].Likes = append(r.Posts[i].Likes, like.UserID)
			}
			return nil
		}
	}
	return model.ErrPostNotFound
}
