package amqp

import (
	"context"
	"encoding/json"
	"k071123/internal/services/notification_service/domain/models"
	"log"
	"time"
)

type EmailSender interface {
	Send(email *models.Email) error
}

type EmailConsumer struct {
	pub  *Publisher
	smtp EmailSender
}

func NewEmailConsumer(pub *Publisher, smtp EmailSender) *EmailConsumer {
	return &EmailConsumer{
		pub:  pub,
		smtp: smtp,
	}
}

func (c *EmailConsumer) Start(ctx context.Context) error {
	ch := c.pub.ch

	deliveries, err := ch.Consume(
		c.pub.queueName,
		"",
		false, // autoAck
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("EmailConsumer stopped")
				return
			case msg, ok := <-deliveries:
				if !ok {
					log.Println("Deliveries channel closed")
					time.Sleep(time.Second)
					continue
				}

				var email models.Email
				if err := json.Unmarshal(msg.Body, &email); err != nil {
					log.Println("Failed to unmarshal email:", err)
					msg.Nack(false, false)
					continue
				}

				if err := c.smtp.Send(&email); err != nil {
					log.Println("Failed to send email:", err)
					msg.Nack(false, true)
					continue
				}

				msg.Ack(false)
				log.Printf("Email sent to %s\n", email.To)
			}
		}
	}()

	return nil
}
