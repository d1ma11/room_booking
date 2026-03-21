package mocks

import (
	"test-backend-1-d1ma11/internal/entity"
	"time"

	"github.com/stretchr/testify/mock"
)

type BookingRepositoryMock struct {
	mock.Mock
}

func (m *BookingRepositoryMock) Create(booking *entity.Booking) error {
	args := m.Called(booking)
	return args.Error(0)
}

func (m *BookingRepositoryMock) Save(booking *entity.Booking) error {
	args := m.Called(booking)
	return args.Error(0)
}

func (m *BookingRepositoryMock) GetById(booking *entity.Booking, bookingId string) error {
	args := m.Called(booking, bookingId)
	return args.Error(0)
}

func (m *BookingRepositoryMock) GetBySlotId(booking *entity.Booking, slotId string) error {
	args := m.Called(booking, slotId)
	return args.Error(0)
}

func (m *BookingRepositoryMock) GetByUserID(bookings *[]entity.Booking, userID string) error {
	args := m.Called(bookings, userID)
	return args.Error(0)
}

func (m *BookingRepositoryMock) GetAllByUserId(bookings *[]entity.Booking, userId string) error {
	args := m.Called(bookings, userId)
	return args.Error(0)
}

func (m *BookingRepositoryMock) GetAllWithPagination(bookings *[]entity.Booking, page, pageSize int, total *int64) error {
	args := m.Called(bookings, page, pageSize, total)
	return args.Error(0)
}

func (m *BookingRepositoryMock) GetBookedSlotIDsByRoomIDAndDate(roomID string, date time.Time, bookedSlotIDs *[]string) error {
	args := m.Called(roomID, date, bookedSlotIDs)
	return args.Error(0)
}
