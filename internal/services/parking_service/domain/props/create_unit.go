package props

import (
	"errors"
	"github.com/google/uuid"
	"k071123/internal/services/parking_service/domain/models"
	"k071123/internal/services/parking_service/domain/models/unit_statuses"
)

type CreateUniqReq struct {
	Status        unit_statuses.UnitStatus        `json:"status"`
	NetworkStatus unit_statuses.NetworkUnitStatus `json:"network_status"`
	Direction     unit_statuses.UnitDirection     `json:"direction"`
	Code          *string                         `json:"code"` // 234654
	QrLink        *string                         `json:"qr_link"`
	ParkingUUID   *uuid.UUID                      `json:"parking_uuid"`
}

type CreateUnitResp struct {
	Unit *models.Unit `json:"unit"`
}

func (r CreateUniqReq) Validate() error {
	if r.Code != nil {
		if len(*r.Code) != 6 {
			return errors.New("code length must be 6")
		}
	}
	return nil
}
