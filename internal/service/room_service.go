package service

import (
	"test-backend-1-d1ma11/internal/entity"
	"test-backend-1-d1ma11/internal/repository"
)

type RoomServiceImpl struct {
	roomRepo repository.RoomRepository
}

func NewRoomService(roomRepo repository.RoomRepository) *RoomServiceImpl {
	return &RoomServiceImpl{roomRepo: roomRepo}
}

func (s *RoomServiceImpl) Create(room *entity.Room) error {
	if err := s.roomRepo.Create(room); err != nil {
		return newInternalError("failed to create room")
	}
	return nil
}

func (s *RoomServiceImpl) List(rooms *[]entity.Room) error {
	if err := s.roomRepo.GetAll(rooms); err != nil {
		return newInternalError("failed to fetch rooms")
	}
	return nil
}
