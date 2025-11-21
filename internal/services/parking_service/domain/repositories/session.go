package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"k071123/internal/services/parking_service/domain/models"
	"time"
)

type SessionUpdates interface {
	SetParkingUUID(value uuid.UUID) SessionUpdates
	SetUnitUUID(value uuid.UUID) SessionUpdates
	SetCarUUID(value uuid.UUID) SessionUpdates
	SetUserUUID(value uuid.UUID) SessionUpdates
	SetFinishAt(finishAt *time.Time) SessionUpdates

	HaveUpdates() bool
}

type SessionRepository interface {
	Filter() SessionFilter
	Updates() SessionUpdates

	Add(model *models.Session) error
	GetByUUID(uuid string) (*models.Session, error)
	WhereFilter(filter SessionFilter) ([]models.Session, error)
	Update(tx *gorm.DB, uuid uuid.UUID, updates SessionUpdates) error
}

type SessionFilter interface {
	SetStatuses([]string) SessionFilter
	SetCarUUIDs([]string) SessionFilter
	SetUUIDs([]string) SessionFilter
}
