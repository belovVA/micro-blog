package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"micro-blog/internal/model"
	"micro-blog/internal/service"
	"micro-blog/internal/service/mocks"
)

func TestUserService_Authenticate(t *testing.T) {
	tests := []struct {
		name         string
		inputUser    *model.User
		mockSetup    func(repo *mocks.UserRepository)
		expectedUser *model.User
		expectErr    bool
	}{
		{
			name:      "existing user",
			inputUser: &model.User{Name: "vova"},
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("GetUserByName", "vova").
					Return(&model.User{Name: "vova"}, nil).
					Once()
			},
			expectedUser: &model.User{Name: "vova"},
			expectErr:    false,
		},
		{
			name:      "new user",
			inputUser: &model.User{Name: "new_user"},
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("GetUserByName", "new_user").
					Return(nil, errors.New("not found")).
					Once()
				repo.On("CreateUser", &model.User{Name: "new_user"}).
					Return(&model.User{Name: "new_user"}, nil).
					Once()
			},
			expectedUser: &model.User{Name: "new_user"},
			expectErr:    false,
		},
		{
			name:      "create user fails",
			inputUser: &model.User{Name: "fail_user"},
			mockSetup: func(repo *mocks.UserRepository) {
				repo.On("GetUserByName", "fail_user").
					Return(nil, errors.New("not found")).
					Once()
				repo.On("CreateUser", &model.User{Name: "fail_user"}).
					Return(nil, errors.New("create failed")).
					Once()
			},
			expectedUser: nil,
			expectErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.UserRepository)
			tt.mockSetup(mockRepo)

			svc := service.NewUserService(mockRepo)

			user, err := svc.Authenticate(context.Background(), tt.inputUser)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser.Name, user.Name)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func BenchmarkUserService_Authenticate(b *testing.B) {
	mockRepo := new(mocks.UserRepository)

	mockRepo.On("GetUserByName", "bench_user").
		Return(&model.User{Name: "bench_user"}, nil).
		Maybe()

	svc := service.NewUserService(mockRepo)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.Authenticate(context.Background(), &model.User{Name: "bench_user"})
	}
}

func TestUserService_Authenticate_Race(t *testing.T) {
	mockRepo := new(mocks.UserRepository)

	mockRepo.On("GetUserByName", "race_user").
		Return(&model.User{Name: "race_user"}, nil).
		Maybe()

	svc := service.NewUserService(mockRepo)

	const goroutines = 50
	errCh := make(chan error, goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			_, err := svc.Authenticate(context.Background(), &model.User{Name: "race_user"})
			errCh <- err
		}()
	}

	for i := 0; i < goroutines; i++ {
		if err := <-errCh; err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}

	mockRepo.AssertExpectations(t)
}
