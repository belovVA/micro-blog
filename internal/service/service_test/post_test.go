package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"micro-blog/internal/model"
	"micro-blog/internal/service"
	mockpost "micro-blog/internal/service/mocks"
	mockuser "micro-blog/internal/service/mocks"
)

func TestPostService_CreatePost(t *testing.T) {
	type fields struct {
		userRepo *mockuser.UserRepository
		postRepo *mockpost.PostRepository
	}

	tests := []struct {
		name       string
		setupMocks func(f fields, post *model.Post)
		post       *model.Post
		wantErr    bool
	}{
		{
			name: "1) User does not exist",
			post: &model.Post{AuthorID: uuid.New(), Text: "Text"},
			setupMocks: func(f fields, post *model.Post) {
				f.userRepo.On("GetUserById", post.AuthorID).Return(nil, model.ErrUserNotFound)
			},
			wantErr: true,
		},
		{
			name: "2) User exists, post created",
			post: &model.Post{AuthorID: uuid.New(), Text: "Hello"},
			setupMocks: func(f fields, post *model.Post) {
				user := &model.User{ID: post.AuthorID, Name: "Alice"}
				f.userRepo.On("GetUserById", post.AuthorID).Return(user, nil)
				f.postRepo.On("CreatePost", post).Return(post, nil)
			},
			wantErr: false,
		},
		{
			name: "3) User exists, post creation fails",
			post: &model.Post{AuthorID: uuid.New(), Text: "Oops"},
			setupMocks: func(f fields, post *model.Post) {
				user := &model.User{ID: post.AuthorID, Name: "Bob"}
				f.userRepo.On("GetUserById", post.AuthorID).Return(user, nil)
				f.postRepo.On("CreatePost", post).Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := mockuser.NewUserRepository(t)
			postRepo := mockpost.NewPostRepository(t)
			tt.setupMocks(fields{userRepo, postRepo}, tt.post)

			s := service.NewPostService(postRepo, userRepo)
			got, err := s.CreatePost(context.Background(), tt.post)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.post, got)
			}
		})
	}
}

func TestPostService_GetListPost(t *testing.T) {
	tests := []struct {
		name       string
		mockReturn []*model.Post
		mockError  error
		wantErr    bool
	}{
		{
			name:       "4) Empty list of posts",
			mockReturn: []*model.Post{},
			mockError:  nil,
			wantErr:    false,
		},
		{
			name: "5) One post in list",
			mockReturn: []*model.Post{
				{
					ID:       uuid.New(),
					AuthorID: uuid.New(),
					Text:     "Hello world",
				},
			},
			mockError: nil,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			postRepo := mockpost.NewPostRepository(t)
			userRepo := mockuser.NewUserRepository(t) // Не используется, но требуется в конструкторе

			postRepo.On("GetListPost").Return(tt.mockReturn, tt.mockError)

			s := service.NewPostService(postRepo, userRepo)
			got, err := s.GetListPost(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockReturn, got)
			}
		})
	}
}
