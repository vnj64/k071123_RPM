package amqp

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"k071123/internal/services/notification_service/domain/models"
)

type Publisher struct {
	ch        *amqp.Channel
	queueName string
}

func NewPublisher(ch *amqp.Channel, queue string) (*Publisher, error) {
	_, err := ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return &Publisher{ch: ch, queueName: queue}, nil
}

func (p *Publisher) Publish(ctx context.Context, msg any) error {
	body, _ := json.Marshal(msg)
	return p.ch.PublishWithContext(ctx, "", p.queueName, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

type AMQPService struct {
	pub *Publisher
}

func NewAMQPService(pub *Publisher) *AMQPService {
	return &AMQPService{pub: pub}
}

// TODO: необходимо бы обобщить
func (s *AMQPService) SendEmail(ctx context.Context, email models.Email) error {
	return s.pub.Publish(ctx, email)
}

func (s *AMQPService) Publisher() *Publisher {
	return s.pub
}
