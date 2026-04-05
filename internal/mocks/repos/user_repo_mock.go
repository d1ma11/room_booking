package repos

import (
	"test-backend-1-d1ma11/internal/entity"

	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) GetByEmail(user *entity.User, email string) error {
	args := m.Called(user, email)
	return args.Error(0)
}

func (m *UserRepositoryMock) Create(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}
