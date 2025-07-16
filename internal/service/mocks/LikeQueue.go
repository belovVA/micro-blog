package mocks

import (
	"github.com/stretchr/testify/mock"
	"micro-blog/internal/model"
)

type MockLikeQueue struct {
	mock.Mock
}

func (m *MockLikeQueue) Enqueue(like *model.Like) {
	m.Called(like)
}
