package repositories

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"k071123/internal/services/parking_service/domain"
	"k071123/internal/services/parking_service/domain/repositories"
	"reflect"
	"sync"
)

type connection struct {
	db *gorm.DB

	mu    sync.Mutex
	cache map[reflect.Type]interface{}
}

type txConnection struct {
	*connection
	tx *gorm.DB
}

func NewConnection(user, password, host, port, database string) (domain.Connection, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", user, password, host, port, database)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	return &connection{
		db:    db,
		cache: make(map[reflect.Type]interface{}),
	}, nil
}

// ______ ТРАНЗАКЦИИ ______

func (c *connection) Begin() (domain.TransactionalConnection, error) {
	tx := c.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &txConnection{
		connection: &connection{
			db:    tx,
			cache: make(map[reflect.Type]interface{}),
		},
		tx: tx,
	}, nil
}

func (t *txConnection) Commit() error {
	return t.tx.Commit().Error
}

func (t *txConnection) Rollback() error {
	return t.tx.Rollback().Error
}

func (c *connection) get(repoPtr interface{}, factory func(*gorm.DB) interface{}) interface{} {
	t := reflect.TypeOf(repoPtr).Elem()

	c.mu.Lock()
	defer c.mu.Unlock()

	if r, ok := c.cache[t]; ok {
		return r
	}

	instance := factory(c.db)
	c.cache[t] = instance
	return instance
}

func (c *connection) DB() *gorm.DB {
	return c.db
}

func (c *connection) CarRepository() repositories.CarRepository {
	return c.get((*repositories.CarRepository)(nil),
		func(db *gorm.DB) interface{} {
			return NewCarRepository(db)
		},
	).(repositories.CarRepository)
}

func (c *connection) ParkingRepository() repositories.ParkingRepository {
	return c.get((*repositories.ParkingRepository)(nil),
		func(db *gorm.DB) interface{} { return NewParkingRepository(db) },
	).(repositories.ParkingRepository)
}

func (c *connection) TariffRepository() repositories.TariffRepository {
	return c.get((*repositories.TariffRepository)(nil),
		func(db *gorm.DB) interface{} { return NewTariffRepository(db) },
	).(repositories.TariffRepository)
}

func (c *connection) UnitRepository() repositories.UnitRepository {
	return c.get((*repositories.UnitRepository)(nil),
		func(db *gorm.DB) interface{} { return NewUnitRepository(db) },
	).(repositories.UnitRepository)
}

func (c *connection) SessionRepository() repositories.SessionRepository {
	return c.get((*repositories.SessionRepository)(nil),
		func(db *gorm.DB) interface{} { return NewSessionRepository(db) },
	).(repositories.SessionRepository)
}
