package repositories

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"k071123/internal/services/parking_service/domain"
	"k071123/internal/services/parking_service/domain/repositories"
)

type connection struct {
	db *gorm.DB

	carRepository     repositories.CarRepository
	parkingRepository repositories.ParkingRepository
	tariffRepository  repositories.TariffRepository
	unitRepository    repositories.UnitRepository
	sessionRepository repositories.SessionRepository
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
		db:                db,
		carRepository:     NewCarRepository(db),
		parkingRepository: NewParkingRepository(db),
		tariffRepository:  NewTariffRepository(db),
		unitRepository:    NewUnitRepository(db),
		sessionRepository: NewSessionRepository(db),
	}, nil
}

func (c *connection) DB() *gorm.DB {
	return c.db
}

func (c *connection) CarRepository() repositories.CarRepository {
	return c.carRepository
}

func (c *connection) ParkingRepository() repositories.ParkingRepository {
	return c.parkingRepository
}

func (c *connection) TariffRepository() repositories.TariffRepository {
	return c.tariffRepository
}

func (c *connection) UnitRepository() repositories.UnitRepository {
	return c.unitRepository
}

func (c *connection) SessionRepository() repositories.SessionRepository {
	return c.sessionRepository
}
