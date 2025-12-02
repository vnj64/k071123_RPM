package repositories

import "k071123/internal/services/order_service/domain/models"

type CardFilter interface {
	SetUUIDs(values []string) CardFilter
	SetLast4(value string) CardFilter
	SetPaymentSystem(value string) CardFilter
	SetUserUUIDs(values []string) CardFilter
	SetIsPreferred(value bool) CardFilter
	SetIsActive(value bool) CardFilter
}

type CardRepository interface {
	Filter() CardFilter

	Insert(model *models.Card) error
	GetByUUID(uuid string) (*models.Card, error)
	GetAllCards(userUUID string) ([]models.Card, error)
	Delete(uuid string) error
	ChangePreferredCard(userUUID string, newCard models.Card) error
	Save(card *models.Card) error
	GetByUserUUID(userUUID string) (*models.Card, error)
	WhereFilter(filter CardFilter) ([]models.Card, error)
}
