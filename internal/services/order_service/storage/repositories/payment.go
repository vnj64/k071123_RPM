package repositories

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"k071123/internal/services/order_service/domain/models"
	"k071123/internal/services/order_service/domain/models/payment_methods"
	"k071123/internal/services/order_service/domain/models/payment_statuses"
	"k071123/internal/services/order_service/domain/repositories"
)

// ________ REPOSITORY ________
type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) repositories.PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Insert(model *models.Payment) error {
	return r.db.Create(model).Error
}

func (r *paymentRepository) GetByTransactionID(id string) (*models.Payment, error) {
	var payment models.Payment
	if err := r.db.Where("transaction_id = ?", id).First(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) Updates() repositories.PaymentUpdates {
	return &paymentUpdates{}
}

func (r *paymentRepository) Pagination() repositories.Pagination {
	return &sqlPagination{}
}

func (r *paymentRepository) Filter() repositories.PaymentFilter {
	return &paymentFilter{}
}

func (r *paymentRepository) Update(uuid uuid.UUID, updates repositories.PaymentUpdates) error {
	u, ok := updates.(*paymentUpdates)
	if !ok {
		return errors.New("invalid payment updates")
	}

	return r.db.Model(&models.Payment{UUID: uuid}).Updates(u.toMap()).Error
}

func (r *paymentRepository) GetPendingPayments() ([]models.Payment, error) {
	var payments []models.Payment
	if err := r.db.Where("status = ?", "pending").Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *paymentRepository) WhereFilter(filter repositories.PaymentFilter, pagination repositories.Pagination) ([]models.Payment, error) {
	var result []models.Payment

	f, ok := filter.(*paymentFilter)
	if !ok {
		return nil, errors.New("invalid filter")
	}

	query := r.db.Model(&models.Payment{})

	query = f.query(query)
	if pagination != nil {
		p, ok := pagination.(*sqlPagination)
		if !ok {
			return nil, errors.New("invalid pagination")
		}
		query = p.query(query)
	}

	// by default
	query = query.Order("created_at desc")
	if err := query.Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

// ________ UPDATES ________
type paymentUpdates struct {
	sessionUUID   *uuid.UUID
	paymentMethod *payment_methods.PaymentMethod
	status        *payment_statuses.PaymentStatuses
	transactionId *string
	amount        *float64
	platformFee   *float64
	description   *string
}

func (u *paymentUpdates) SetSessionUUID(value uuid.UUID) repositories.PaymentUpdates {
	u.sessionUUID = &value
	return u
}

func (u *paymentUpdates) SetPaymentMethod(value payment_methods.PaymentMethod) repositories.PaymentUpdates {
	u.paymentMethod = &value
	return u
}

func (u *paymentUpdates) SetStatus(value payment_statuses.PaymentStatuses) repositories.PaymentUpdates {
	u.status = &value
	return u
}

func (u *paymentUpdates) SetTransactionId(value string) repositories.PaymentUpdates {
	u.transactionId = &value
	return u
}

func (u *paymentUpdates) SetAmount(value float64) repositories.PaymentUpdates {
	u.amount = &value
	return u
}

func (u *paymentUpdates) SetPlatformFee(value float64) repositories.PaymentUpdates {
	u.platformFee = &value
	return u
}

func (u *paymentUpdates) SetDescription(value string) repositories.PaymentUpdates {
	u.description = &value
	return u
}

func (u *paymentUpdates) HaveUpdates() bool {
	return len(u.toMap()) > 0
}

func (u *paymentUpdates) toMap() map[string]interface{} {
	out := make(map[string]interface{})

	if u.sessionUUID != nil {
		out["session_uuid"] = *u.sessionUUID
	}
	if u.paymentMethod != nil {
		out["payment_method"] = *u.paymentMethod
	}
	if u.status != nil {
		out["status"] = *u.status
	}
	if u.transactionId != nil {
		out["transaction_id"] = *u.transactionId
	}
	if u.amount != nil {
		out["amount"] = *u.amount
	}
	if u.platformFee != nil {
		out["platform_fee"] = *u.platformFee
	}
	if u.description != nil {
		out["description"] = *u.description
	}
	return out
}

// ______ FILTER ______

type paymentFilter struct {
	sessionUUIDs []string
	statuses     []string
}

func (f *paymentFilter) SetSessionUUIDs(value []string) repositories.PaymentFilter {
	f.sessionUUIDs = value
	return f
}

func (f *paymentFilter) SetStatuses(value []string) repositories.PaymentFilter {
	f.statuses = value
	return f
}

func (f *paymentFilter) query(tx *gorm.DB) *gorm.DB {
	if len(f.sessionUUIDs) > 0 {
		tx = tx.Where("session_uuid IN (?)", f.sessionUUIDs)
	}

	if len(f.statuses) > 0 {
		tx = tx.Where("status IN (?)", f.statuses)
	}

	return tx
}
