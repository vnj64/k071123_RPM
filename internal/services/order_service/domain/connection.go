package domain

import (
	"gorm.io/gorm"
	"k071123/internal/services/order_service/domain/repositories"
)

type Connection interface {
	DB() *gorm.DB
	Card() repositories.CardRepository
	Payment() repositories.PaymentRepository
	VerifyToken() repositories.VerifyTokenRepository
}
