package models

import (
	"github.com/google/uuid"
	"k071123/internal/services/parking_service/domain/models/parking_statuses"
	"k071123/pkg/timestamps"
)

type Parking struct {
	Timestamps  timestamps.Timestamps          `gorm:"embedded;" json:"-" swaggerignore:"true"`
	UUID        uuid.UUID                      `gorm:"primaryKey" json:"uuid"`
	TariffUUID  uuid.UUID                      `json:"tariff_uuid"`
	Name        string                         `json:"name"`
	Address     string                         `json:"address"`
	Latitude    string                         `json:"latitude"`
	Longitude   string                         `json:"longitude"`
	WorkingTime []ParkingSchedule              `gorm:"foreignKey:ParkingUUID" json:"working_time"`
	TotalPlaces uint                           `json:"total_places"`
	Status      parking_statuses.ParkingStatus `json:"status"`
}

func (p Parking) TableName() string {
	return "parkings"
}
