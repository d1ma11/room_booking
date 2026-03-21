package mocks

import (
	"test-backend-1-d1ma11/internal/entity"

	"github.com/stretchr/testify/mock"
)

type ScheduleRepositoryMock struct {
	mock.Mock
}

func (m *ScheduleRepositoryMock) Create(schedule *entity.Schedule) error {
	args := m.Called(schedule)
	return args.Error(0)
}

func (m *ScheduleRepositoryMock) GetByRoomId(schedule *entity.Schedule, roomId string) error {
	args := m.Called(schedule, roomId)
	return args.Error(0)
}
