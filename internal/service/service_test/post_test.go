package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"micro-blog/internal/model"
	"micro-blog/internal/service"
	mockpost "micro-blog/internal/service/mocks"
	mockqueue "micro-blog/internal/service/mocks"
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

func TestPostService_LikePost(t *testing.T) {
	type args struct {
		like *model.Like
	}

	userID := uuid.New()
	postID := uuid.New()

	tests := []struct {
		name           string
		args           args
		mockUser       func(*mockuser.UserRepository)
		mockPost       func(*mockpost.PostRepository)
		mockLikeQueue  func(*mockqueue.MockLikeQueue)
		expectedErrMsg error
		runTwice       bool
		secondErr      error
	}{
		{
			name: "user not found",
			args: args{like: &model.Like{PostID: postID, UserID: userID}},
			mockUser: func(ur *mockuser.UserRepository) {
				ur.On("GetUserById", userID).Return(nil, model.ErrUserNotFound)
			},
			mockPost:       func(pr *mockpost.PostRepository) {}, // не нужен тут
			mockLikeQueue:  func(lq *mockqueue.MockLikeQueue) {}, // не нужен тут
			expectedErrMsg: model.ErrUserNotFound,
		},
		{
			name: "send enqueue like",
			args: args{like: &model.Like{PostID: postID, UserID: userID}},
			mockUser: func(ur *mockuser.UserRepository) {
				ur.On("GetUserById", userID).Return(&model.User{ID: userID, Name: "Alice"}, nil)
			},
			mockPost: func(pr *mockpost.PostRepository) {
			},
			mockLikeQueue: func(lq *mockqueue.MockLikeQueue) {
				lq.On("Enqueue", mock.Anything).Once()

			},
			expectedErrMsg: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := mockuser.NewUserRepository(t)
			postRepo := mockpost.NewPostRepository(t)
			likeQueue := new(mockqueue.MockLikeQueue)

			tt.mockUser(userRepo)
			tt.mockPost(postRepo)
			tt.mockLikeQueue(likeQueue)

			ps := service.NewPostService(postRepo, userRepo)
			ps.AttachLikeQueue(likeQueue)

			err := ps.LikePost(context.Background(), tt.args.like)
			assert.Equal(t, tt.expectedErrMsg, err)

			userRepo.AssertExpectations(t)
			postRepo.AssertExpectations(t)
			likeQueue.AssertExpectations(t)
		})
	}
}

func BenchmarkPostService_LikePost(b *testing.B) {
	userID := uuid.New()
	postID := uuid.New()

	userRepo := mockuser.NewUserRepository(b)
	postRepo := mockpost.NewPostRepository(b)
	likeQueue := new(mockqueue.MockLikeQueue)

	userRepo.On("GetUserById", mock.Anything).Return(&model.User{ID: userID, Name: "BenchUser"}, nil)
	likeQueue.On("Enqueue", mock.Anything).Return()

	service := service.NewPostService(postRepo, userRepo)
	service.AttachLikeQueue(likeQueue)

	like := &model.Like{PostID: postID, UserID: userID}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.LikePost(context.Background(), like)
	}
}

func TestPostService_LikePost_Concurrent(t *testing.T) {
	userID := uuid.New()
	postID := uuid.New()

	userRepo := mockuser.NewUserRepository(t)
	postRepo := mockpost.NewPostRepository(t)
	likeQueue := new(mockqueue.MockLikeQueue)

	userRepo.On("GetUserById", userID).Return(&model.User{ID: userID, Name: "ConcurrentUser"}, nil)
	likeQueue.On("Enqueue", mock.Anything).Return()

	service := service.NewPostService(postRepo, userRepo)
	service.AttachLikeQueue(likeQueue)

	like := &model.Like{PostID: postID, UserID: userID}

	t.Run("parallel likes", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 100; i++ {
			go func() {
				err := service.LikePost(context.Background(), like)
				assert.NoError(t, err)
			}()
		}
	})
}
