package entity

import "time"

type Slot struct {
	ID      string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	RoomID  string    `gorm:"type:uuid;not null;index:idx_slots_room_date" json:"roomId"`
	StartAt time.Time `gorm:"type:timestamptz;not null;index:idx_slots_room_date" json:"start"`
	EndAt   time.Time `gorm:"type:timestamptz;not null" json:"end"`
}

func (Slot) TableName() string {
	return "slots"
}
