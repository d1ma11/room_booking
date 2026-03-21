package service

import (
	"errors"
	"test-backend-1-d1ma11/internal/entity"
	"test-backend-1-d1ma11/internal/repository"

	"gorm.io/gorm"
)

type ScheduleServiceImpl struct {
	roomRepo     repository.RoomRepository
	scheduleRepo repository.ScheduleRepository
}

func NewScheduleService(roomRepo repository.RoomRepository, scheduleRepo repository.ScheduleRepository) *ScheduleServiceImpl {
	return &ScheduleServiceImpl{roomRepo: roomRepo, scheduleRepo: scheduleRepo}
}

func (s *ScheduleServiceImpl) Create(schedule *entity.Schedule, roomId string) error {
	var room entity.Room
	if err := s.roomRepo.GetById(&room, roomId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NewError(ErrorType.RoomNotFound, "room not found")
		}
		return newInternalError("failed to check room")
	}

	if err := s.scheduleRepo.GetByRoomId(&entity.Schedule{}, roomId); err == nil {
		return NewError(ErrorType.ScheduleExists, "schedule for this room already exists and cant be changed")
	}

	if err := s.scheduleRepo.Create(schedule); err != nil {
		return newInternalError("failed to create schedule")
	}
	return nil

}
