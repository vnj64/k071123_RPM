package props

import (
	"errors"
	"k071123/internal/services/user_service/domain/models"
	"k071123/internal/services/user_service/domain/models/user_status"
	"time"
)

type CreateUserReq struct {
	FirstName   *string            `json:"first_name"`
	SecondName  *string            `json:"second_name"`
	BirthDate   *time.Time         `json:"birth_date"`
	Status      user_status.Status `json:"status"`
	PhoneNumber string             `json:"phone_number"`
}

type CreateUserResp struct {
	User *models.User `json:"user"`
}

func (r CreateUserReq) Validate() error {
	if r.PhoneNumber == "" {
		return errors.New("phone number cannot be empty")
	}
	return nil
}
