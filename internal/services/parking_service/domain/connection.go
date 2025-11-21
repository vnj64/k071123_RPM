package domain

import (
	"gorm.io/gorm"
	"k071123/internal/services/parking_service/domain/repositories"
)

type Connection interface {
	DB() *gorm.DB
	CarRepository() repositories.CarRepository
	ParkingRepository() repositories.ParkingRepository
	TariffRepository() repositories.TariffRepository
	UnitRepository() repositories.UnitRepository
	SessionRepository() repositories.SessionRepository
	Begin() (TransactionalConnection, error)
}

type TransactionalConnection interface {
	Connection

	Commit() error
	Rollback() error
}
