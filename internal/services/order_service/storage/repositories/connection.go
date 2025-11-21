package repositories

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"k071123/internal/services/order_service/domain"
	"k071123/internal/services/order_service/domain/repositories"
)

type connection struct {
	db *gorm.DB

	cardRepository        repositories.CardRepository
	paymentRepository     repositories.PaymentRepository
	verifyTokenRepository repositories.VerifyTokenRepository
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
		db:                    db,
		cardRepository:        NewCardRepository(db),
		paymentRepository:     NewPaymentRepository(db),
		verifyTokenRepository: NewVerifyTokenRepository(db),
	}, nil
}

func (c *connection) DB() *gorm.DB {
	return c.db
}

func (c *connection) Card() repositories.CardRepository {
	return c.cardRepository
}

func (c *connection) Payment() repositories.PaymentRepository {
	return c.paymentRepository
}

func (c *connection) VerifyToken() repositories.VerifyTokenRepository {
	return c.verifyTokenRepository
}
