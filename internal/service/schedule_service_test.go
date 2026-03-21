package service

import (
	"errors"
	"test-backend-1-d1ma11/internal/entity"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func init() {
	scheduleService = NewScheduleService(mockRoomRepo, mockScheduleRepo)
}

func TestScheduleService_Create_Success(t *testing.T) {
	// Arrange
	roomId := "test-room-uuid"
	schedule := &entity.Schedule{
		RoomID:     roomId,
		DaysOfWeek: pq.Int32Array{1, 2, 3},
		StartTime:  "09:00:00",
		EndTime:    "18:00:00",
	}

	mockRoomRepo.On("GetById", mock.AnythingOfType("*entity.Room"), roomId).Return(nil).Once()
	mockScheduleRepo.On("GetByRoomId", mock.AnythingOfType("*entity.Schedule"), roomId).Return(gorm.ErrRecordNotFound).Once()
	mockScheduleRepo.On("Create", schedule).Return(nil).Once()

	// Act
	err := scheduleService.Create(schedule, roomId)

	// Assert
	assert.NoError(t, err)
	mockRoomRepo.AssertExpectations(t)
	mockScheduleRepo.AssertExpectations(t)
}

func TestScheduleService_Create_RoomNotFound(t *testing.T) {
	// Arrange
	roomId := "non-existent-room"
	schedule := &entity.Schedule{RoomID: roomId}

	mockRoomRepo.On("GetById", mock.AnythingOfType("*entity.Room"), roomId).Return(gorm.ErrRecordNotFound).Once()

	// Act
	err := scheduleService.Create(schedule, roomId)

	// Assert
	assert.Error(t, err)
	assert.IsType(t, &ErrorResponse{}, err)
	appErr := err.(*ErrorResponse)
	assert.Equal(t, ErrorType.RoomNotFound, appErr.Code)
	assert.Equal(t, "room not found", appErr.Message)
	mockRoomRepo.AssertExpectations(t)
	mockScheduleRepo.AssertNotCalled(t, "GetByRoomId")
	mockScheduleRepo.AssertNotCalled(t, "Create")
}

func TestScheduleService_Create_ScheduleAlreadyExists(t *testing.T) {
	// Arrange
	roomId := "test-room-uuid"
	schedule := &entity.Schedule{RoomID: roomId}

	mockRoomRepo.On("GetById", mock.AnythingOfType("*entity.Room"), roomId).Return(nil).Once()
	mockScheduleRepo.On("GetByRoomId", mock.AnythingOfType("*entity.Schedule"), roomId).Return(nil).Once()

	// Act
	err := scheduleService.Create(schedule, roomId)

	// Assert
	assert.Error(t, err)
	assert.IsType(t, &ErrorResponse{}, err)
	appErr := err.(*ErrorResponse)
	assert.Equal(t, ErrorType.ScheduleExists, appErr.Code)
	assert.Equal(t, "schedule for this room already exists and cant be changed", appErr.Message)
	mockRoomRepo.AssertExpectations(t)
	mockScheduleRepo.AssertExpectations(t)
	mockScheduleRepo.AssertNotCalled(t, "Create")
}

func TestScheduleService_Create_RoomRepoError(t *testing.T) {
	// Arrange
	roomId := "test-room-uuid"
	schedule := &entity.Schedule{RoomID: roomId}
	dbErr := errors.New("connection failed")

	mockRoomRepo.On("GetById", mock.AnythingOfType("*entity.Room"), roomId).Return(dbErr).Once()

	// Act
	err := scheduleService.Create(schedule, roomId)

	// Assert
	assert.Error(t, err)
	assert.IsType(t, &InternalErrorResponse{}, err)
	assert.Equal(t, ErrorType.InternalError, err.(*InternalErrorResponse).Code)
	mockRoomRepo.AssertExpectations(t)
	mockScheduleRepo.AssertNotCalled(t, "GetByRoomId")
}

func TestScheduleService_Create_ScheduleRepoCreateError(t *testing.T) {
	// Arrange
	roomId := "test-room-uuid"
	schedule := &entity.Schedule{RoomID: roomId}
	dbErr := errors.New("insert failed")

	mockRoomRepo.On("GetById", mock.AnythingOfType("*entity.Room"), roomId).Return(nil).Once()
	mockScheduleRepo.On("GetByRoomId", mock.AnythingOfType("*entity.Schedule"), roomId).Return(gorm.ErrRecordNotFound).Once()
	mockScheduleRepo.On("Create", schedule).Return(dbErr).Once()

	// Act
	err := scheduleService.Create(schedule, roomId)

	// Assert
	assert.Error(t, err)
	assert.IsType(t, &InternalErrorResponse{}, err)
	assert.Equal(t, ErrorType.InternalError, err.(*InternalErrorResponse).Code)
	mockRoomRepo.AssertExpectations(t)
	mockScheduleRepo.AssertExpectations(t)
}
