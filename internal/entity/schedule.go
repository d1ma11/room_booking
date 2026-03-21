package entity

import "github.com/lib/pq"

type Schedule struct {
	ID         string        `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	RoomID     string        `gorm:"type:uuid;not null;uniqueIndex" json:"roomId"`
	DaysOfWeek pq.Int32Array `gorm:"type:int[];not null" json:"daysOfWeek"`
	StartTime  string        `gorm:"type:time;not null" json:"startTime"`
	EndTime    string        `gorm:"type:time;not null" json:"endTime"`
}

func (Schedule) TableName() string {
	return "schedules"
}
