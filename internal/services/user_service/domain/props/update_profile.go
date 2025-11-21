package props

import (
	"errors"
	"k071123/internal/services/user_service/domain/models"
	"k071123/pkg/validation"
	"time"
)

type UpdateProfileReq struct {
	UserUUID    string     `json:"-"`
	FirstName   *string    `json:"first_name"`
	SecondName  *string    `json:"second_name"`
	BirthDate   *time.Time `json:"birth_date"`
	PhoneNumber *string    `json:"phone_number"`
}

func (r *UpdateProfileReq) Validate() error {
	if r.UserUUID == "" {
		return errors.New("user_uuid is required")
	}
	if r.PhoneNumber != nil {
		if err := validation.ValidatePhoneNumber(*r.PhoneNumber); err != nil {
			return err
		}
	}
	return nil
}

type UpdateProfileResp struct {
	User *models.User `json:"user"`
}
