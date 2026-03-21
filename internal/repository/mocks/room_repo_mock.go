package mocks

import (
	"test-backend-1-d1ma11/internal/entity"

	"github.com/stretchr/testify/mock"
)

type RoomRepositoryMock struct {
	mock.Mock
}

func (m *RoomRepositoryMock) Create(room *entity.Room) error {
	args := m.Called(room)
	return args.Error(0)
}

func (m *RoomRepositoryMock) GetAll(rooms *[]entity.Room) error {
	args := m.Called(rooms)
	return args.Error(0)
}

func (m *RoomRepositoryMock) GetById(room *entity.Room, roomId string) error {
	args := m.Called(room, roomId)
	return args.Error(0)
}
