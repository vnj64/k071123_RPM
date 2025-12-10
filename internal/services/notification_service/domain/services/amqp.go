package services

import (
	"context"
	"k071123/internal/services/notification_service/domain/models"
	"k071123/internal/services/notification_service/services/amqp"
)

type AMQP interface {
	SendEmail(ctx context.Context, msg models.Email) error
	Publisher() *amqp.Publisher
}
