package service

import (
	"test-backend-1-d1ma11/configs"
	"test-backend-1-d1ma11/internal/entity"
	"test-backend-1-d1ma11/internal/repository"
	"test-backend-1-d1ma11/internal/service/auth"
	"time"
)

type RoomService interface {
	Create(room *entity.Room) error
	List(rooms *[]entity.Room) error
}

type ScheduleService interface {
	Create(schedule *entity.Schedule, roomId string) error
}

type SlotService interface {
	ListAvailableSlots(roomID string, date time.Time) ([]entity.Slot, error)
	GenerateSlotsForDate(roomID string, date time.Time) ([]entity.Slot, error)
}

type BookingService interface {
	Create(booking *entity.Booking, createConferenceLink bool) error
	ListByUserId(bookings *[]entity.Booking, userId string) error
	CancelBooking(booking *entity.Booking, bookingId, userId string) error
	ListByAdminId(bookings *[]entity.Booking, page, pageSize int, total *int64) error
}

type ConferenceService interface {
	CreateConferenceLink(bookingID string) (string, error)
}

type UserService interface {
	DummyLogin(role string) (string, error)
	Login(email, password string) (string, error)
	Register(email, password, role string) (entity.User, error)
}

type JwtService interface {
	ParseToken(tokenStr, secret string) (*auth.Claims, error)
	GenerateToken(userId, role, secret string, ttl time.Duration) (string, error)
}

type Services struct {
	RoomService     RoomService
	ScheduleService ScheduleService
	SlotService     SlotService
	BookingService  BookingService
	UserService     UserService
	JwtService      JwtService
}

type ServicesDependencies struct {
	Repos             *repository.Repositories
	ConferenceService ConferenceService
	Config            *configs.Config
}

func NewServices(deps ServicesDependencies) *Services {
	jwtSvc := auth.NewJwtService()
	return &Services{
		RoomService:     NewRoomService(deps.Repos.RoomRepository),
		ScheduleService: NewScheduleService(deps.Repos.RoomRepository, deps.Repos.ScheduleRepository),
		SlotService:     NewSlotServiceImpl(deps.Repos.BookingRepository, deps.Repos.SlotRepository, deps.Repos.ScheduleRepository, deps.Repos.RoomRepository),
		BookingService:  NewBookingService(deps.Repos.BookingRepository, deps.Repos.SlotRepository, deps.ConferenceService),
		UserService:     NewUserServiceImpl(deps.Config, deps.Repos.UserRepository, jwtSvc),
		JwtService:      jwtSvc,
	}
}
