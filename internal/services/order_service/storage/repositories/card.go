package repositories

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"k071123/internal/services/order_service/domain/models"
	"k071123/internal/services/order_service/domain/repositories"
	"log"
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
			log.Printf("unable to find user card: [%s]", userUUID)
		}
	}
	if err := r.db.Model(&card).Where("uuid = ?", card.UUID.String()).Update("is_preferred", false).Error; err != nil {
		return err
	}
	if err := r.db.Create(&newCard).Error; err != nil {
		return err
	}

	return nil
}

func (r *cardRepository) GetByUserUUID(userUUID string) (*models.Card, error) {
	var card models.Card
	if err := r.db.Where("user_uuid = ?", userUUID).Order("created_at DESC").First(&card).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &card, nil
}

func (r *cardRepository) Save(card *models.Card) error {
	return r.db.Save(card).Error
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

// ______ FILTER ______

type cardFilter struct {
	uuids         []string
	last4         *string
	paymentSystem *string
	userUUIDs     []string
	isPreferred   *bool
	isActive      *bool
}

func (r *cardRepository) Filter() repositories.CardFilter {
	return &cardFilter{}
}

// TODO: добавить в WhereFilter аргумент OrderBy
func (r *cardRepository) WhereFilter(filter repositories.CardFilter) ([]models.Card, error) {
	var cards []models.Card
	f, ok := filter.(*cardFilter)
	if !ok {
		return []models.Card{}, errors.New("wrong filter type")
	}
	query := f.query(r.db)
	if err := r.db.Where(query).Find(&cards).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.Card{}, nil
		}
		return nil, err
	}
	return cards, nil
}

func (f *cardFilter) query(tx *gorm.DB) *gorm.DB {
	if len(f.uuids) > 0 {
		tx.Where("uuid IN (?)", f.uuids)
	}
	if len(f.userUUIDs) > 0 {
		tx.Where("user_uuid IN (?)", f.userUUIDs)
	}
	if f.last4 != nil {
		tx.Where("last4 = ?", *f.last4)
	}
	if f.paymentSystem != nil {
		tx.Where("payment_system = ?", *f.paymentSystem)
	}
	if f.isActive != nil {
		tx.Where("is_active = ?", *f.isActive)
	}
	if f.isPreferred != nil {
		tx.Where("is_preferred = ?", *f.isPreferred)
	}
	tx.Order("created_at DESC")
	
	return tx
}

func (f *cardFilter) SetUUIDs(uuids []string) repositories.CardFilter {
	f.uuids = uuids
	return f
}

func (f *cardFilter) SetLast4(last4 string) repositories.CardFilter {
	f.last4 = &last4
	return f
}

func (f *cardFilter) SetPaymentSystem(paymentSystem string) repositories.CardFilter {
	f.paymentSystem = &paymentSystem
	return f
}

func (f *cardFilter) SetUserUUIDs(userUUIDs []string) repositories.CardFilter {
	f.userUUIDs = userUUIDs
	return f
}

func (f *cardFilter) SetIsActive(isActive bool) repositories.CardFilter {
	f.isActive = &isActive
	return f
}

func (f *cardFilter) SetIsPreferred(isPreferred bool) repositories.CardFilter {
	f.isPreferred = &isPreferred
	return f
}
