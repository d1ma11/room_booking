package entity

import "time"

type User struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:uuid_generated_v4()" json:"id"`
	Email     string    `gorm:"type:text;uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"type:text;not null" json:"-"`
	Role      string    `gorm:"type:text;not null; check: role IN ('admin','user')" json:"role"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}

func (User) TableName() string {
	return "users"
}
