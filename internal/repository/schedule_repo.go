package repository

import (
	"test-backend-1-d1ma11/internal/entity"

	"gorm.io/gorm"
)

type ScheduleRepositoryImpl struct {
	db *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) *ScheduleRepositoryImpl {
	return &ScheduleRepositoryImpl{db: db}
}

func (r *ScheduleRepositoryImpl) GetByRoomId(schedule *entity.Schedule, roomId string) error {
	if err := r.db.Where("room_id = ?", roomId).First(&schedule).Error; err != nil {
		return err
	}
	return nil
}

func (r *ScheduleRepositoryImpl) Create(schedule *entity.Schedule) error {
	if err := r.db.Create(&schedule).Error; err != nil {
		return err
	}
	return nil
}
