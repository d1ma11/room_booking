package repository

import (
	"test-backend-1-d1ma11/internal/entity"

	"gorm.io/gorm"
)

type RoomRepositoryImpl struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *RoomRepositoryImpl {
	return &RoomRepositoryImpl{db: db}
}

func (r *RoomRepositoryImpl) Create(room *entity.Room) error {
	if err := r.db.Create(&room).Error; err != nil {
		return err
	}
	return nil
}

func (r *RoomRepositoryImpl) GetAll(rooms *[]entity.Room) error {
	if err := r.db.Find(&rooms).Error; err != nil {
		return err
	}
	return nil
}

func (r *RoomRepositoryImpl) GetById(room *entity.Room, roomId string) error {
	if err := r.db.First(&room, "id = ?", roomId).Error; err != nil {
		return err
	}
	return nil
}
