package repositories

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"k071123/internal/services/user_service/domain"
	"k071123/internal/services/user_service/domain/repositories"
)

type connection struct {
	db                   *gorm.DB
	userRepo             repositories.User
	verificationCodeRepo repositories.VerificationCode
}

func NewConnection(user, password, host, port, database string) (domain.Connection, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", user, password, host, port, database)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	return &connection{
		db:                   db,
		userRepo:             NewUserRepository(db),
		verificationCodeRepo: NewVerificationCodeRepository(db),
	}, nil
}

func (c *connection) User() repositories.User {
	return c.userRepo
}

func (c *connection) VerificationCode() repositories.VerificationCode {
	return c.verificationCodeRepo
}

func (c *connection) DB() *gorm.DB {
	return c.db
}
