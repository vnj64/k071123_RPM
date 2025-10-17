package repositories

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"k071123/internal/services/notification_service/domain"
)

type connection struct {
	db *gorm.DB
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
		db: db,
	}, nil
}

func (c *connection) DB() *gorm.DB {
	return c.db
}
