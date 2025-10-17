package core

import (
	"k071123/internal/services/notification_service/domain"
	domainServices "k071123/internal/services/notification_service/domain/services"
	"k071123/internal/services/notification_service/services/config"
	"k071123/internal/services/notification_service/services/smtp"
	"k071123/internal/services/notification_service/storage/repositories"
)

type Ctx struct {
	services   domain.Services
	connection domain.Connection
}

type svs struct {
	config domainServices.Config
	smtp   domainServices.Smtp
}

func (s *svs) Config() domainServices.Config {
	return s.config
}

func (s *svs) Smtp() domainServices.Smtp {
	return s.smtp
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

	return &Ctx{
		services: &svs{
			config: cfg,
			smtp:   smtpClient,
		},
		connection: sqlConnection,
	}
}
