package service

import (
	"errors"
	"fmt"
	log "log/slog"
	"test-backend-1-d1ma11/internal/entity"
	"test-backend-1-d1ma11/internal/repository"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SlotServiceImpl struct {
	bookingRepo  repository.BookingRepository
	slotRepo     repository.SlotRepository
	scheduleRepo repository.ScheduleRepository
	roomRepo     repository.RoomRepository
}

func NewSlotServiceImpl(
	bookingRepo repository.BookingRepository,
	slotRepo repository.SlotRepository,
	scheduleRepo repository.ScheduleRepository,
	roomRepo repository.RoomRepository) *SlotServiceImpl {
	return &SlotServiceImpl{bookingRepo: bookingRepo, slotRepo: slotRepo, scheduleRepo: scheduleRepo, roomRepo: roomRepo}
}

func (s *SlotServiceImpl) GenerateSlotsForDate(roomId string, date time.Time) ([]entity.Slot, error) {
	var room entity.Room
	if err := s.roomRepo.GetById(&room, roomId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info(fmt.Sprintf("Room with room_id=%v was not found. Returning zero slots", roomId))
			return []entity.Slot{}, NewError(ErrorType.RoomNotFound, "room not found")
		}
		log.Info(fmt.Sprintf("Internal error during getting room by room_id=%v. Error: %v", roomId, err))
		return []entity.Slot{}, newInternalError("failed to check room")
	}

	var existing []entity.Slot
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, time.UTC)

	if err := s.slotRepo.GetAll(roomId, startOfDay, endOfDay, &existing); err != nil {
		log.Info(fmt.Sprintf("Get an error while getting slots from DB. Error: %v", err))
		return nil, err
	}
	if len(existing) > 0 {
		log.Info(fmt.Sprintf("Slots have been already generated, returning existing slots. Slots: %v", existing))
		return existing, nil
	}

	var schedule entity.Schedule
	if err := s.scheduleRepo.GetByRoomId(&schedule, roomId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info("Schedule was not found. Returning zero slots")
			return []entity.Slot{}, nil
		}
		log.Info(fmt.Sprintf("Internal error during getting room by room_idid=%v. Error: %v", roomId, err))
		return nil, err
	}

	dayOfWeek := int(date.Weekday())
	if dayOfWeek == 0 {
		dayOfWeek = 7
	}

	isValidDay := false
	for _, d := range schedule.DaysOfWeek {
		if d == int32(dayOfWeek) {
			isValidDay = true
			break
		}
	}
	if !isValidDay {
		log.Info(fmt.Sprintf("This room is not working on this day: %v", date.Day()))
		return []entity.Slot{}, nil
	}

	startTime, err := time.Parse("15:04:05", schedule.StartTime)
	if err != nil {
		log.Info(fmt.Sprintf("Invalid startTime: %v", err))
		return nil, fmt.Errorf("invalid startTime: %w", err)
	}
	endTime, err := time.Parse("15:04:05", schedule.EndTime)
	if err != nil {
		log.Info(fmt.Sprintf("Invalid endTime: %v", err))
		return nil, fmt.Errorf("invalid endTime: %w", err)
	}

	startLimit := time.Date(startOfDay.Year(), startOfDay.Month(), startOfDay.Day(),
		startTime.Hour(), startTime.Minute(), 0, 0, time.UTC)
	endLimit := time.Date(startOfDay.Year(), startOfDay.Month(), startOfDay.Day(),
		endTime.Hour(), endTime.Minute(), 0, 0, time.UTC)

	var slots []entity.Slot
	current := startLimit

	for !current.Add(30 * time.Minute).After(endLimit) {
		slotEnd := current.Add(30 * time.Minute)
		slot := entity.Slot{
			ID:      uuid.New().String(),
			RoomID:  roomId,
			StartAt: current,
			EndAt:   slotEnd,
		}
		slots = append(slots, slot)
		current = slotEnd
	}

	if len(slots) > 0 {
		log.Info("Saving generated slots in table 'slots'")
		if err = s.slotRepo.CreateAll(&slots); err != nil {
			log.Info(fmt.Sprintf("There is an error during saving slots. Error: %v", err))
			_ = s.slotRepo.GetAllByRoomIdInDate(roomId, startOfDay, endOfDay, &slots)
		}
	}

	return slots, nil
}

func (s *SlotServiceImpl) ListAvailableSlots(roomID string, date time.Time) ([]entity.Slot, error) {
	slots, err := s.GenerateSlotsForDate(roomID, date)
	if err != nil {
		return nil, err
	}

	var bookedSlotIDs []string
	if err = s.bookingRepo.GetBookedSlotIDsByRoomIDAndDate(roomID, date, &bookedSlotIDs); err != nil {
		return nil, err
	}

	bookedMap := make(map[string]bool)
	for _, id := range bookedSlotIDs {
		bookedMap[id] = true
	}

	var available []entity.Slot
	for _, slot := range slots {
		if !bookedMap[slot.ID] {
			available = append(available, slot)
		}
	}

	return available, nil
}
