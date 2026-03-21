package service

import "test-backend-1-d1ma11/internal/repository/mocks"

var (
	mockBookingRepo  *mocks.BookingRepositoryMock
	mockSlotRepo     *mocks.SlotRepositoryMock
	mockScheduleRepo *mocks.ScheduleRepositoryMock
	mockRoomRepo     *mocks.RoomRepositoryMock
	slotService      SlotService
	roomService      RoomService
	scheduleService  ScheduleService
)

func init() {
	mockBookingRepo = new(mocks.BookingRepositoryMock)
	mockSlotRepo = new(mocks.SlotRepositoryMock)
	mockScheduleRepo = new(mocks.ScheduleRepositoryMock)
	mockRoomRepo = new(mocks.RoomRepositoryMock)
}
