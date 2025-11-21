package repositories

import (
	"github.com/google/uuid"
	"k071123/internal/services/order_service/domain/models"
	"k071123/internal/services/order_service/domain/models/payment_methods"
	"k071123/internal/services/order_service/domain/models/payment_statuses"
)

type PaymentFilter interface {
	SetSessionUUIDs(value []string) PaymentFilter
	SetStatuses(value []string) PaymentFilter
}

type PaymentUpdates interface {
	SetSessionUUID(value uuid.UUID) PaymentUpdates
	SetPaymentMethod(value payment_methods.PaymentMethod) PaymentUpdates
	SetStatus(value payment_statuses.PaymentStatuses) PaymentUpdates
	SetTransactionId(value string) PaymentUpdates
	SetAmount(value float64) PaymentUpdates
	SetPlatformFee(value float64) PaymentUpdates
	SetDescription(value string) PaymentUpdates

	HaveUpdates() bool
}

type PaymentRepository interface {
	Updates() PaymentUpdates
	Filter() PaymentFilter
	Pagination() Pagination

	Insert(model *models.Payment) error
	Update(uuid uuid.UUID, updates PaymentUpdates) error
	GetByTransactionID(id string) (*models.Payment, error)
	GetPendingPayments() ([]models.Payment, error)
	WhereFilter(filter PaymentFilter, pagination Pagination) ([]models.Payment, error)
}
