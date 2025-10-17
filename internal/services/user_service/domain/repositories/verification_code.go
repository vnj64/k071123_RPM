package repositories

import "k071123/internal/services/user_service/domain/models"

type VerificationCode interface {
	GetByCode(code string) (*models.VerificationCode, error)
	Add(model *models.VerificationCode) error
	GetByUUID(uuid string) (*models.VerificationCode, error)
	GetByEmail(email string) (*models.VerificationCode, error)
}
