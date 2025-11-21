package repositories

import "k071123/internal/services/order_service/domain/models"

type CardRepository interface {
	Insert(model *models.Card) error
	GetByUUID(uuid string) (*models.Card, error)
	GetAllCards(userUUID string) ([]models.Card, error)
	Delete(uuid string) error
	ChangePreferredCard(userUUID string, newCard models.Card) error
}
