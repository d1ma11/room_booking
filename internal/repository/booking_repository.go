package repository

import (
	"test-backend-1-d1ma11/internal/entity"
	"time"

	"gorm.io/gorm"
)

type BookingRepositoryImpl struct {
	DB *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepositoryImpl {
	return &BookingRepositoryImpl{DB: db}
}

func (r *BookingRepositoryImpl) GetBookedSlotIDsByRoomIDAndDate(roomID string, date time.Time, bookedSlotIDs *[]string) error {
	if err := r.DB.Table("bookings").
		Select("slot_id").
		Joins("JOIN slots ON bookings.slot_id = slots.id").
		Where("slots.room_id = ? AND slots.start_at >= ? AND slots.start_at < ? AND bookings.status = ?",
			roomID,
			time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC),
			time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, time.UTC),
			"active").
		Pluck("slot_id", &bookedSlotIDs).Error; err != nil {
		return err
	}
	return nil
}

func (r *BookingRepositoryImpl) Create(booking *entity.Booking) error {
	if err := r.DB.Create(&booking).Error; err != nil {
		return err
	}
	return nil
}

func (r *BookingRepositoryImpl) GetAllByUserId(bookings *[]entity.Booking, userId string) error {
	if err := r.DB.Find(&bookings).Where("user_id = ?", userId).Error; err != nil {
		return err
	}
	return nil
}

func (r *BookingRepositoryImpl) Save(booking *entity.Booking) error {
	if err := r.DB.Save(booking).Error; err != nil {
		return err
	}
	return nil
}

func (r *BookingRepositoryImpl) GetBySlotId(booking *entity.Booking, slotId string) error {
	if err := r.DB.Where("slot_id = ? AND status = 'active'", slotId).First(&booking).Error; err != nil {
		return err
	}
	return nil
}

func (r *BookingRepositoryImpl) GetById(booking *entity.Booking, bookingId string) error {
	if err := r.DB.First(&booking, "id = ?", bookingId).Error; err != nil {
		return err
	}
	return nil
}

func (r *BookingRepositoryImpl) GetByUserID(bookings *[]entity.Booking, userID string) error {
	if err := r.DB.
		Preload("Slot").
		Joins("JOIN slots ON bookings.slot_id = slots.id").
		Where("bookings.user_id = ? AND slots.start_at >= ?", userID, time.Now().UTC()).
		Find(bookings).Error; err != nil {
		return err
	}
	return nil
}
func (r *BookingRepositoryImpl) GetAllWithPagination(bookings *[]entity.Booking, page, pageSize int, total *int64) error {
	if err := r.DB.Model(&entity.Booking{}).Count(total).Error; err != nil {
		return err
	}

	offset := (page - 1) * pageSize
	if err := r.DB.Preload("Slot").
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(bookings).Error; err != nil {
		return err
	}
	return nil
}
