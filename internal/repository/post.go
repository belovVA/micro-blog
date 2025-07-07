package repository

import (
	"github.com/google/uuid"
	"micro-blog/internal/model"
)

const (
	initPostsCapacity = 100
)

type PostRepo struct {
	Posts []*model.Post
}

func NewPostRepo() *PostRepo {
	return &PostRepo{
		Posts: make([]*model.Post, 0, initPostsCapacity),
	}
}

func (r *PostRepo) CreatePost(post *model.Post) (*model.Post, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	post.ID = id
	r.Posts = append(r.Posts, post)
	return post, nil
}

func (r *PostRepo) GetListPost() ([]*model.Post, error) {
	return r.Posts, nil
}

func (r *PostRepo) LikePost(like *model.Like) error {
	for i, post := range r.Posts {
		if post.ID == like.PostID {
			r.Posts[i].Likes = append(r.Posts[i].Likes, like.UserID)
			return nil
		}
	}
	return model.ErrPostNotFound
}
