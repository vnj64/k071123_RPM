package models

import (
	"github.com/google/uuid"
)

type ParkingSchedule struct {
	UUID        uuid.UUID `gorm:"gorm:primaryKey" json:"uuid"`
	DaysOfWeek  []int     `gorm:"type:integer[]" json:"days_of_week"`
	OpenTime    string    `json:"open_time"`
	CloseTime   string    `json:"close_time"`
	ParkingUUID uuid.UUID `gorm:"gorm:constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"parking_uuid"`
}

func (p ParkingSchedule) TableName() string {
	return "parking_schedule"
}

// Monday 10.00-20.00
// Tuesday 24hrs
