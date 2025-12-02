package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ParkingSchedule struct {
	UUID        uuid.UUID     `gorm:"primaryKey" json:"uuid"`
	DaysOfWeek  pq.Int64Array `gorm:"type:int[]" json:"days_of_week"`
	OpenTime    string        `json:"open_time"`
	CloseTime   string        `json:"close_time"`
	ParkingUUID uuid.UUID     `gorm:"gorm:constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"parking_uuid"`
}

func (p ParkingSchedule) TableName() string {
	return "parking_schedule"
}

// Monday 10.00-20.00
// Tuesday 24hrs
