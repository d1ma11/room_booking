package service

import (
	"errors"
	"test-backend-1-d1ma11/internal/entity"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func init() {
	slotService = NewSlotServiceImpl(mockBookingRepo, mockSlotRepo, mockScheduleRepo, mockRoomRepo)
}

func TestSlotService_GenerateSlotsForDate_AlreadyExist(t *testing.T) {
	// Arrange
	roomID := "room-uuid"
	testDate := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
	startOfDay := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
	endOfDay := time.Date(2024, 6, 10, 23, 59, 59, 999999999, time.UTC)
	existingSlots := []entity.Slot{
		{ID: "slot-1", RoomID: roomID, StartAt: startOfDay, EndAt: startOfDay.Add(30 * time.Minute)},
	}

	mockRoomRepo.On("GetById", mock.AnythingOfType("*entity.Room"), roomID).Return(nil).Once()
	mockSlotRepo.On("GetAll", roomID, startOfDay, endOfDay, mock.AnythingOfType("*[]entity.Slot")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(3).(*[]entity.Slot)
		*arg = existingSlots
	}).Once()

	// Act
	slots, err := slotService.GenerateSlotsForDate(roomID, testDate)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, slots, 1)
	assert.Equal(t, "slot-1", slots[0].ID)
	mockRoomRepo.AssertExpectations(t)
	mockSlotRepo.AssertExpectations(t)
	mockScheduleRepo.AssertNotCalled(t, "GetByRoomId")
}

func TestSlotService_GenerateSlotsForDate_NoSchedule(t *testing.T) {
	// Arrange
	roomID := "room-uuid"
	testDate := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
	startOfDay := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
	endOfDay := time.Date(2024, 6, 10, 23, 59, 59, 999999999, time.UTC)

	mockRoomRepo.On("GetById", mock.AnythingOfType("*entity.Room"), roomID).Return(nil).Once()
	mockSlotRepo.On("GetAll", roomID, startOfDay, endOfDay, mock.AnythingOfType("*[]entity.Slot")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(3).(*[]entity.Slot)
		*arg = []entity.Slot{}
	}).Once()
	mockScheduleRepo.On("GetByRoomId", mock.AnythingOfType("*entity.Schedule"), roomID).Return(gorm.ErrRecordNotFound).Once()

	// Act
	slots, err := slotService.GenerateSlotsForDate(roomID, testDate)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, slots)
	mockRoomRepo.AssertExpectations(t)
	mockSlotRepo.AssertExpectations(t)
	mockScheduleRepo.AssertExpectations(t)
}

func TestSlotService_GenerateSlotsForDate_InvalidDayOfWeek(t *testing.T) {
	// Arrange
	roomID := "room-uuid"
	testDate := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
	startOfDay := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
	endOfDay := time.Date(2024, 6, 10, 23, 59, 59, 999999999, time.UTC)
	schedule := entity.Schedule{
		RoomID:     roomID,
		DaysOfWeek: pq.Int32Array{2},
		StartTime:  "09:00:00",
		EndTime:    "18:00:00",
	}

	mockRoomRepo.On("GetById", mock.AnythingOfType("*entity.Room"), roomID).Return(nil).Once()
	mockSlotRepo.On("GetAll", roomID, startOfDay, endOfDay, mock.AnythingOfType("*[]entity.Slot")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(3).(*[]entity.Slot)
		*arg = []entity.Slot{}
	}).Once()

	mockScheduleRepo.On("GetByRoomId", mock.AnythingOfType("*entity.Schedule"), roomID).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*entity.Schedule)
		*arg = schedule
	}).Once()

	// Act
	slots, err := slotService.GenerateSlotsForDate(roomID, testDate)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, slots)
	mockRoomRepo.AssertExpectations(t)
	mockSlotRepo.AssertExpectations(t)
	mockScheduleRepo.AssertExpectations(t)
}

func TestSlotService_GenerateSlotsForDate_Success(t *testing.T) {
	// Arrange
	roomID := "room-uuid"
	testDate := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
	startOfDay := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
	endOfDay := time.Date(2024, 6, 10, 23, 59, 59, 999999999, time.UTC)
	schedule := entity.Schedule{
		RoomID:     roomID,
		DaysOfWeek: pq.Int32Array{1},
		StartTime:  "09:00:00",
		EndTime:    "10:00:00",
	}

	mockRoomRepo.On("GetById", mock.AnythingOfType("*entity.Room"), roomID).Return(nil).Once()
	mockSlotRepo.On("GetAll", roomID, startOfDay, endOfDay, mock.AnythingOfType("*[]entity.Slot")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(3).(*[]entity.Slot)
		*arg = []entity.Slot{}
	}).Once()
	mockScheduleRepo.On("GetByRoomId", mock.AnythingOfType("*entity.Schedule"), roomID).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*entity.Schedule)
		*arg = schedule
	}).Once()
	mockSlotRepo.On("CreateAll", mock.AnythingOfType("*[]entity.Slot")).Return(nil).Once()

	// Act
	slots, err := slotService.GenerateSlotsForDate(roomID, testDate)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, slots, 2)
	assert.Equal(t, time.Date(2024, 6, 10, 9, 0, 0, 0, time.UTC), slots[0].StartAt)
	assert.Equal(t, time.Date(2024, 6, 10, 9, 30, 0, 0, time.UTC), slots[0].EndAt)
	assert.Equal(t, time.Date(2024, 6, 10, 9, 30, 0, 0, time.UTC), slots[1].StartAt)
	assert.Equal(t, time.Date(2024, 6, 10, 10, 0, 0, 0, time.UTC), slots[1].EndAt)
	mockRoomRepo.AssertExpectations(t)
	mockSlotRepo.AssertExpectations(t)
	mockScheduleRepo.AssertExpectations(t)
}

func TestSlotService_GenerateSlotsForDate_CreateAllError(t *testing.T) {
	// Arrange
	roomID := "room-uuid"
	testDate := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
	startOfDay := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
	endOfDay := time.Date(2024, 6, 10, 23, 59, 59, 999999999, time.UTC)
	dbErr := errors.New("insert failed")
	schedule := entity.Schedule{
		RoomID:     roomID,
		DaysOfWeek: pq.Int32Array{1},
		StartTime:  "09:00:00",
		EndTime:    "10:00:00",
	}

	mockRoomRepo.On("GetById", mock.AnythingOfType("*entity.Room"), roomID).Return(nil).Once()
	mockSlotRepo.On("GetAll", roomID, startOfDay, endOfDay, mock.AnythingOfType("*[]entity.Slot")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(3).(*[]entity.Slot)
		*arg = []entity.Slot{}
	}).Once()
	mockScheduleRepo.On("GetByRoomId", mock.AnythingOfType("*entity.Schedule"), roomID).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*entity.Schedule)
		*arg = schedule
	}).Once()
	mockSlotRepo.On("CreateAll", mock.AnythingOfType("*[]entity.Slot")).Return(dbErr).Once()
	mockSlotRepo.On("GetAllByRoomIdInDate", roomID, startOfDay, endOfDay, mock.AnythingOfType("*[]entity.Slot")).Return(nil).Once()

	// Act
	slots, err := slotService.GenerateSlotsForDate(roomID, testDate)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, slots)
	mockRoomRepo.AssertExpectations(t)
	mockSlotRepo.AssertExpectations(t)
	mockScheduleRepo.AssertExpectations(t)
}

func TestSlotService_GenerateSlotsForDate_RoomNotFound(t *testing.T) {
	// Arrange
	roomID := "non-existent-room"
	testDate := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)

	mockRoomRepo.On("GetById", mock.AnythingOfType("*entity.Room"), roomID).Return(gorm.ErrRecordNotFound).Once()

	// Act
	slots, err := slotService.GenerateSlotsForDate(roomID, testDate)

	// Assert
	assert.Error(t, err)
	assert.IsType(t, &ErrorResponse{}, err)
	appErr := err.(*ErrorResponse)
	assert.Equal(t, ErrorType.RoomNotFound, appErr.Code)
	assert.Empty(t, slots)
	mockRoomRepo.AssertExpectations(t)
	mockSlotRepo.AssertNotCalled(t, "GetAll")
	mockScheduleRepo.AssertNotCalled(t, "GetByRoomId")
}

func TestSlotService_ListAvailableSlots_BookingRepoError(t *testing.T) {
	// Arrange
	roomID := "room-uuid"
	testDate := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
	startOfDay := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
	endOfDay := time.Date(2024, 6, 10, 23, 59, 59, 999999999, time.UTC)
	dbErr := errors.New("booking repo error")
	schedule := entity.Schedule{
		RoomID:     roomID,
		DaysOfWeek: pq.Int32Array{1},
		StartTime:  "09:00:00",
		EndTime:    "10:00:00",
	}

	mockRoomRepo.On("GetById", mock.AnythingOfType("*entity.Room"), roomID).Return(nil).Once()
	mockSlotRepo.On("GetAll", roomID, startOfDay, endOfDay, mock.AnythingOfType("*[]entity.Slot")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(3).(*[]entity.Slot)
		*arg = []entity.Slot{}
	}).Once()
	mockScheduleRepo.On("GetByRoomId", mock.AnythingOfType("*entity.Schedule"), roomID).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*entity.Schedule)
		*arg = schedule
	}).Once()
	mockSlotRepo.On("CreateAll", mock.AnythingOfType("*[]entity.Slot")).Return(nil).Once()
	mockBookingRepo.On("GetBookedSlotIDsByRoomIDAndDate", roomID, testDate, mock.AnythingOfType("*[]string")).Return(dbErr).Once()

	// Act
	slots, err := slotService.ListAvailableSlots(roomID, testDate)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, slots)
	mockRoomRepo.AssertExpectations(t)
	mockBookingRepo.AssertExpectations(t)
}
