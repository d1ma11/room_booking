package entity

import "time"

type Room struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name        string    `gorm:"type:text;not null" json:"name"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`
	Capacity    *int      `gorm:"type:int" json:"capacity,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
}

func (Room) TableName() string {
	return "rooms"
}
