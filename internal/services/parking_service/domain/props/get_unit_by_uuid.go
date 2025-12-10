package props

import (
	"github.com/google/uuid"
	"k071123/internal/services/parking_service/domain/models"
)

type GetUnitByUUID struct {
	UUID uuid.UUID `json:"uuid9"`
}

type GetUnitByUUIDResp struct {
	Unit *models.Unit `json:"unit"`
}
