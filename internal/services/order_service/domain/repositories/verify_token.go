package repositories

import "k071123/internal/services/order_service/domain/models"

type VerifyTokenRepository interface {
	Insert(model *models.VerifyTokens) error
	GetLastByUserUUID(userUuid string) (*models.VerifyTokens, error)
	Save(model *models.VerifyTokens) error
}
