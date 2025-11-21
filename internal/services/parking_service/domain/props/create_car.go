package props

import (
	"errors"
	"k071123/internal/services/parking_service/domain/models"
)

type CreateCarReq struct {
	GosNumber string `json:"gos_number"`
	UserUUID  string `json:"user_uuid"`
}

func (r *CreateCarReq) Validate() error {
	if r.GosNumber == "" {
		return errors.New("gos_number is required")
	}
	return nil
}

type CreateCarResp struct {
	Car *models.Car `json:"car"`
}
