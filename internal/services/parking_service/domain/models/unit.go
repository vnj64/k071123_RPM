package models

import (
	"github.com/google/uuid"
	"k071123/internal/services/parking_service/domain/models/unit_statuses"
)

type Unit struct {
	UUID          uuid.UUID                       `gorm:"type:uuid;primaryKey" json:"uuid"`
	Status        unit_statuses.UnitStatus        `json:"status"`
	NetworkStatus unit_statuses.NetworkUnitStatus `json:"network_status"`
	Direction     unit_statuses.UnitDirection     `json:"direction"`
	Code          *string                         `json:"code"` // 234654
	QrLink        *string                         `json:"qr_link"`
	ParkingUUID   *uuid.UUID                      `json:"parking_uuid"`
}
