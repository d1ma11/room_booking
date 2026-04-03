package repository

import (
	"test-backend-1-d1ma11/internal/entity"
	"time"

	"gorm.io/gorm"
)

type RoomRepository interface {
	Create(room *entity.Room) error
	GetAll(rooms *[]entity.Room) error
	GetById(room *entity.Room, roomId string) error
}

type ScheduleRepository interface {
	Create(schedule *entity.Schedule) error
	GetByRoomId(schedule *entity.Schedule, roomId string) error
}

type SlotRepository interface {
	CreateAll(slots *[]entity.Slot) error
	GetById(slot *entity.Slot, slotId string) error
	GetAll(roomID string, startOfDay, endOfDay time.Time, slots *[]entity.Slot) error
	GetAllByRoomIdInDate(roomId string, from, to time.Time, slots *[]entity.Slot) error
}

type BookingRepository interface {
	Save(booking *entity.Booking) error
	Create(booking *entity.Booking) error
	GetById(booking *entity.Booking, bookingId string) error
	GetBySlotId(booking *entity.Booking, slotId string) error
	GetByUserID(bookings *[]entity.Booking, userID string) error
	GetAllByUserId(bookings *[]entity.Booking, userId string) error
	GetAllWithPagination(bookings *[]entity.Booking, page, pageSize int, total *int64) error
	GetBookedSlotIDsByRoomIDAndDate(roomID string, date time.Time, bookedSlotIDs *[]string) error
}

type UserRepository interface {
	Create(user *entity.User) error
	GetByEmail(user *entity.User, email string) error
}

type Repositories struct {
	RoomRepository     RoomRepository
	ScheduleRepository ScheduleRepository
	SlotRepository     SlotRepository
	BookingRepository  BookingRepository
	UserRepository     UserRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		RoomRepository:     NewRoomRepository(db),
		ScheduleRepository: NewScheduleRepository(db),
		SlotRepository:     NewSlotRepository(db),
		BookingRepository:  NewBookingRepository(db),
		UserRepository:     NewUserRepository(db),
	}
}
