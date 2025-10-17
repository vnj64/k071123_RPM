package domain

import (
	"gorm.io/gorm"
	"k071123/internal/services/user_service/domain/repositories"
)

type Connection interface {
	DB() *gorm.DB
	User() repositories.User
	VerificationCode() repositories.VerificationCode
}
