package entity

import "time"

type Booking struct {
	ID             string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	SlotID         string    `gorm:"type:uuid;not null;uniqueIndex:idx_unique_active_booking,where:status = 'active'" json:"slotId"`
	UserID         string    `gorm:"type:uuid;not null;index" json:"userId"`
	Status         string    `gorm:"type:text;not null;check:status IN ('active','cancelled')" json:"status"`
	ConferenceLink *string   `gorm:"type:text" json:"conferenceLink,omitempty"`
	CreatedAt      time.Time `json:"createdAt,omitempty"`
	Slot           Slot      `gorm:"foreignKey:SlotID" json:"-"`
}

func (Booking) TableName() string {
	return "bookings"
}
