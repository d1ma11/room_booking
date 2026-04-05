package service

import (
	"test-backend-1-d1ma11/configs"
	//"test-backend-1-d1ma11/internal/controller"
	mockRepos "test-backend-1-d1ma11/internal/mocks/repos"
	mockService "test-backend-1-d1ma11/internal/mocks/services"
	"time"
)

var (
	// Repository mocks
	mockUserRepo     *mockRepos.UserRepositoryMock
	mockRoomRepo     *mockRepos.RoomRepositoryMock
	mockSlotRepo     *mockRepos.SlotRepositoryMock
	mockBookingRepo  *mockRepos.BookingRepositoryMock
	mockScheduleRepo *mockRepos.ScheduleRepositoryMock

	// Service mocks
	mockJwtService  *mockService.JwtServiceMock
	mockUserSvc     *mockService.UserServiceMock
	mockRoomSvc     *mockService.RoomServiceMock
	mockSlotSvc     *mockService.SlotServiceMock
	mockBookingSvc  *mockService.BookingServiceMock
	mockScheduleSvc *mockService.ScheduleServiceMock

	// Test configuration
	testConfig = &configs.Config{
		JWT: configs.JWTConfig{
			Secret: "test-secret",
			TTL:    time.Hour},
	}

	// Testing services
	userService     UserService
	roomService     RoomService
	slotService     SlotService
	bookingService  BookingService
	scheduleService ScheduleService

	// Testing controllers
	//userCtrl controller.UserController
	//roomCtrl controller.RoomController
	//slotCtrl controller.SlotController
	//bookingCtrl  controller.BookingController
	//scheduleCtrl controller.ScheduleController
)

func init() {
	mockJwtService = new(mockService.JwtServiceMock)

	mockUserRepo = new(mockRepos.UserRepositoryMock)
	mockRoomRepo = new(mockRepos.RoomRepositoryMock)
	mockSlotRepo = new(mockRepos.SlotRepositoryMock)
	mockBookingRepo = new(mockRepos.BookingRepositoryMock)
	mockScheduleRepo = new(mockRepos.ScheduleRepositoryMock)
}
