package repository

import (
	"test-backend-1-d1ma11/internal/entity"
	"time"

	"gorm.io/gorm"
)

type SlotRepositoryImpl struct {
	db *gorm.DB
}

func NewSlotRepository(db *gorm.DB) *SlotRepositoryImpl {
	return &SlotRepositoryImpl{db: db}
}

func (r *SlotRepositoryImpl) GetAll(roomID string, startOfDay, endOfDay time.Time, slots *[]entity.Slot) error {
	if err := r.db.Where("room_id = ? AND start_at >= ? AND start_at < ?",
		roomID, startOfDay, endOfDay).Find(&slots).Error; err != nil {
		return err
	}
	return nil
}

func (r *SlotRepositoryImpl) CreateAll(slots *[]entity.Slot) error {
	if err := r.db.Create(&slots).Error; err != nil {
		return err
	}
	return nil
}

func (r *SlotRepositoryImpl) GetAllByRoomIdInDate(roomId string, from, to time.Time, slots *[]entity.Slot) error {
	if err := r.db.Where("room_id = ? AND start_at >= ? AND start_at < ?",
		roomId, from, to).Find(&slots).Error; err != nil {
		return err
	}
	return nil
}

func (r *SlotRepositoryImpl) GetById(slot *entity.Slot, slotId string) error {
	if err := r.db.First(&slot, "id = ?", slotId).Error; err != nil {
		return err
	}
	return nil
}
