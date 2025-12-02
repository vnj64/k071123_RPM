package repositories

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"k071123/internal/services/parking_service/domain/models"
	"k071123/internal/services/parking_service/domain/repositories"
	"time"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

// ______ REPOSITORY ______

func (r *SessionRepository) Add(model *models.Session) error {
	return r.db.Create(model).Error
}

func (r *SessionRepository) GetByUUID(uuid string) (*models.Session, error) {
	var car models.Session
	if err := r.db.Where("uuid = ?", uuid).First(&car).Error; err != nil {
		return nil, err
	}
	return &car, nil
}

func (r *SessionRepository) WhereFilter(filter repositories.SessionFilter) ([]models.Session, error) {
	var parkings []models.Session
	f, ok := filter.(*sessionFilter)
	if !ok {
		return []models.Session{}, errors.New("wrong filter type")
	}
	query := f.query(r.db)
	if err := r.db.Where(query).Find(&parkings).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.Session{}, nil
		}
		return nil, err
	}
	return parkings, nil
}

func (r *SessionRepository) Update(uuid uuid.UUID, updates repositories.SessionUpdates) error {
	var session models.Session
	if err := r.db.Where("uuid = ?", uuid).First(&session).Error; err != nil {
		return err
	}
	if err := r.db.Model(&session).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}

func (r *SessionRepository) Filter() repositories.SessionFilter {
	return &sessionFilter{}
}

func (r *SessionRepository) Updates() repositories.SessionUpdates {
	return &sessionUpdates{}
}

// ______ FILTER ______
type sessionFilter struct {
	carUUIDs []string
	statuses []string
	uuids    []string
}

func (f *sessionFilter) SetCarUUIDs(carUUIDs []string) repositories.SessionFilter {
	f.carUUIDs = carUUIDs
	return f
}

func (f *sessionFilter) SetStatuses(statuses []string) repositories.SessionFilter {
	f.statuses = statuses
	return f
}

func (f *sessionFilter) SetUUIDs(uuids []string) repositories.SessionFilter {
	f.uuids = uuids
	return f
}

func (f *sessionFilter) query(tx *gorm.DB) *gorm.DB {
	if len(f.carUUIDs) > 0 {
		tx.Where("car_uuid IN (?)", f.carUUIDs)
	}
	if len(f.statuses) > 0 {
		tx.Where("status IN (?)", f.statuses)
	}
	if len(f.uuids) > 0 {
		tx.Where("uuid IN (?)", f.uuids)
	}
	return tx
}

// ______ UPDATES ______
type sessionUpdates struct {
	parkingUUID *uuid.UUID
	unitUUID    *uuid.UUID
	carUUID     *uuid.UUID
	userUUID    *uuid.UUID
	status      *string
	finishAt    *time.Time
	cost        *float64
}

func (u *sessionUpdates) toMap() map[string]interface{} {
	out := make(map[string]interface{})
	if u.parkingUUID != nil {
		out["parking_uuid"] = u.parkingUUID
	}
	if u.unitUUID != nil {
		out["unit_uuid"] = u.unitUUID
	}
	if u.carUUID != nil {
		out["car_uuid"] = u.carUUID
	}
	if u.userUUID != nil {
		out["user_uuid"] = u.userUUID
	}
	if u.finishAt != nil {
		out["finish_at"] = u.finishAt
	}
	if u.status != nil {
		out["status"] = u.status
	}
	if u.cost != nil {
		out["cost"] = u.cost
	}
	return out
}

func (u *sessionUpdates) SetParkingUUID(parkingUUID uuid.UUID) repositories.SessionUpdates {
	u.parkingUUID = &parkingUUID
	return u
}

func (u *sessionUpdates) SetUnitUUID(unitUUID uuid.UUID) repositories.SessionUpdates {
	u.unitUUID = &unitUUID
	return u
}

func (u *sessionUpdates) SetCarUUID(carUUID uuid.UUID) repositories.SessionUpdates {
	u.carUUID = &carUUID
	return u
}

func (u *sessionUpdates) SetUserUUID(userUUID uuid.UUID) repositories.SessionUpdates {
	u.userUUID = &userUUID
	return u
}

func (u *sessionUpdates) SetFinishAt(finishAt *time.Time) repositories.SessionUpdates {
	u.finishAt = finishAt
	return u
}

func (u *sessionUpdates) SetStatus(status string) repositories.SessionUpdates {
	u.status = &status
	return u
}

func (u *sessionUpdates) SetCost(cost float64) repositories.SessionUpdates {
	u.cost = &cost
	return u
}

func (u *sessionUpdates) HaveUpdates() bool {
	return len(u.toMap()) > 0
}
