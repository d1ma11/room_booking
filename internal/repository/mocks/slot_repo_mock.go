package mocks

import (
	"test-backend-1-d1ma11/internal/entity"
	"time"

	"github.com/stretchr/testify/mock"
)

type SlotRepositoryMock struct {
	mock.Mock
}

func (m *SlotRepositoryMock) CreateAll(slots *[]entity.Slot) error {
	args := m.Called(slots)
	return args.Error(0)
}

func (m *SlotRepositoryMock) GetById(slot *entity.Slot, slotId string) error {
	args := m.Called(slot, slotId)
	return args.Error(0)
}

func (m *SlotRepositoryMock) GetAll(roomID string, startOfDay, endOfDay time.Time, slots *[]entity.Slot) error {
	args := m.Called(roomID, startOfDay, endOfDay, slots)
	return args.Error(0)
}

func (m *SlotRepositoryMock) GetAllByRoomIdInDate(roomId string, from, to time.Time, slots *[]entity.Slot) error {
	args := m.Called(roomId, from, to, slots)
	return args.Error(0)
}
