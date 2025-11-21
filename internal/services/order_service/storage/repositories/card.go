package repositories

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"k071123/internal/services/order_service/domain/models"
	"k071123/internal/services/order_service/domain/repositories"
	"time"
)

type cardRepository struct {
	db *gorm.DB
}

func NewCardRepository(db *gorm.DB) repositories.CardRepository {
	return &cardRepository{db: db}
}

func (r *cardRepository) Insert(model *models.Card) error {
	return r.db.Create(model).Error
}

func (r *cardRepository) GetByUUID(uuid string) (*models.Card, error) {
	var model models.Card
	if err := r.db.Where("uuid = ?", uuid).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &models.Card{}, nil
		}
		return nil, err
	}
	return &model, nil
}

func (r *cardRepository) GetAllCards(userUUID string) ([]models.Card, error) {
	var dbCards []models.Card
	if err := r.db.Where(&models.Card{UserUUID: userUUID}).Find(&dbCards).Error; err != nil {
		return nil, fmt.Errorf("getting cards for user %q in db: %w", userUUID, err)
	}
	cards := make([]models.Card, len(dbCards))
	for i := range dbCards {
		cards[i] = models.Card{
			UUID:          dbCards[i].UUID,
			Last4Digits:   dbCards[i].Last4Digits,
			PaymentSystem: dbCards[i].PaymentSystem,
			UserUUID:      dbCards[i].UserUUID,
		}
	}
	return cards, nil
}

func (r *cardRepository) SetPreferredCard(uuid string) error {
	return nil
}

func (r *cardRepository) ChangePreferredCard(userUUID string, newCard models.Card) error {
	var card models.Card
	if err := r.db.Where("user_uuid = ? AND is_preferred = true", userUUID).First(&card).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := r.db.Create(&newCard).Error; err != nil {
				return err
			}
		}
	}
	if err := r.db.Model(&card).Update("is_preferred", false).Error; err != nil {
		return err
	}
	if err := r.db.Create(&newCard).Error; err != nil {
		return err
	}

	return nil
}

func (r *cardRepository) Delete(uuid string) error {
	var model models.Card
	if err := r.db.Where("uuid = ?", uuid).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("card not found")
		}
		return err
	}

	if err := r.db.Model(&model).Update("deleted_at", gorm.DeletedAt{Time: time.Now(), Valid: true}).Error; err != nil {
		return err
	}

	return nil
}
