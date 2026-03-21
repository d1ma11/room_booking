package service

import (
	"errors"
	"test-backend-1-d1ma11/internal/entity"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	roomService = NewRoomService(mockRoomRepo)
}

func TestRoomService_Create_Success(t *testing.T) {
	// Arrange
	room := &entity.Room{
		Name:        "Test Room",
		Description: stringPtr("Description"),
		Capacity:    intPtr(10),
	}

	mockRoomRepo.On("Create", room).Return(nil).Once()

	// Act
	err := roomService.Create(room)

	// Assert
	assert.NoError(t, err)
	mockRoomRepo.AssertExpectations(t)
}

func TestRoomService_Create_RepositoryError(t *testing.T) {
	// Arrange
	room := &entity.Room{Name: "Test Room"}
	dbErr := errors.New("database error")

	mockRoomRepo.On("Create", room).Return(dbErr).Once()

	// Act
	err := roomService.Create(room)

	// Assert
	assert.Error(t, err)
	assert.IsType(t, &InternalErrorResponse{}, err)
	assert.Equal(t, ErrorType.InternalError, err.(*InternalErrorResponse).Code)
	mockRoomRepo.AssertExpectations(t)
}

func TestRoomService_List_Success(t *testing.T) {
	// Arrange
	var rooms []entity.Room
	mockRoomRepo.On("GetAll", mock.AnythingOfType("*[]entity.Room")).Return(nil).Once()

	// Act
	err := roomService.List(&rooms)

	// Assert
	assert.NoError(t, err)
	mockRoomRepo.AssertExpectations(t)
}

func TestRoomService_List_RepositoryError(t *testing.T) {
	// Arrange
	var rooms []entity.Room
	dbErr := errors.New("database error")

	mockRoomRepo.On("GetAll", mock.AnythingOfType("*[]entity.Room")).Return(dbErr).Once()

	// Act
	err := roomService.List(&rooms)

	// Assert
	assert.Error(t, err)
	assert.IsType(t, &InternalErrorResponse{}, err)
	assert.Equal(t, ErrorType.InternalError, err.(*InternalErrorResponse).Code)
	mockRoomRepo.AssertExpectations(t)
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
