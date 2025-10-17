package services

import "k071123/internal/services/notification_service/domain/models"

type Smtp interface {
	Send(email *models.Email) error
}
