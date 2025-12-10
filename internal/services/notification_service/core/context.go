package core

import (
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"k071123/internal/services/notification_service/domain"
	domainServices "k071123/internal/services/notification_service/domain/services"
	"k071123/internal/services/notification_service/services/amqp"
	"k071123/internal/services/notification_service/services/config"
	"k071123/internal/services/notification_service/services/smtp"
	"k071123/internal/services/notification_service/storage/repositories"
	"k071123/tools/logger"
)

type Ctx struct {
	services   domain.Services
	connection domain.Connection
}

type svs struct {
	config domainServices.Config
	smtp   domainServices.Smtp
	logger *logger.Logger
	amqp   domainServices.AMQP
}

func (s *svs) Config() domainServices.Config {
	return s.config
}

func (s *svs) Smtp() domainServices.Smtp {
	return s.smtp
}

func (s *svs) Logger() *logger.Logger {
	return s.logger
}

func (s *svs) Amqp() domainServices.AMQP {
	return s.amqp
}

func (c *Ctx) Services() domain.Services {
	return c.services
}

func (c *Ctx) Connection() domain.Connection {
	return c.connection
}

func (c *Ctx) Make() domain.Context {
	return &Ctx{
		services:   c.services,
		connection: c.connection,
	}
}

func InitCtx() *Ctx {
	cfg := config.Make()

	loggerCfg := logger.Config{
		Host:     cfg.ElasticHost(),
		Port:     cfg.ElasticPort(),
		Username: cfg.ElasticUsername(),
		Password: cfg.ElasticPassword(),
		Index:    "notification",
		Service:  "notification_service",
	}
	log, err := logger.New(loggerCfg)
	if err != nil {
		panic(err)
	}

	smtpClient := smtp.NewSmtpClient(cfg)
	sqlConnection, err := repositories.NewConnection(
		cfg.PostgresUser(),
		cfg.PostgresPassword(),
		cfg.PostgresHost(),
		cfg.PostgresPort(),
		cfg.PostgresDbName(),
	)
	if err != nil {
		panic("connection isnt success")
	}

	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.AMQPUser(), cfg.AMQPPassword(), cfg.AMQPHost(), cfg.AMQPPort())
	conn, err := amqp091.Dial(url)
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	pub, err := amqp.NewPublisher(ch, "email_notifications")
	if err != nil {
		panic(err)
	}
	amqpSvc := amqp.NewAMQPService(pub)

	return &Ctx{
		services: &svs{
			config: cfg,
			smtp:   smtpClient,
			logger: log,
			amqp:   amqpSvc,
		},
		connection: sqlConnection,
	}
}
