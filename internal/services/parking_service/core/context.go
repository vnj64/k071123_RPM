package core

import (
	"github.com/sirupsen/logrus"
	"k071123/internal/services/parking_service/domain"
	domainServices "k071123/internal/services/parking_service/domain/services"
	"k071123/internal/services/parking_service/services/config"
	"k071123/internal/services/parking_service/storage/repositories"
	"k071123/tools/logger"
)

type Ctx struct {
	services   domain.Services
	connection domain.Connection
}

type svs struct {
	config domainServices.Config
	logger *logrus.Logger
}

func (s *svs) Config() domainServices.Config {
	return s.config
}

func (s *svs) Logger() *logrus.Logger {
	return s.logger
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
		Index:    "parking",
	}
	log, err := logger.New(loggerCfg)
	if err != nil {
		panic(err)
	}

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
			logger: log,
		},
		connection: sqlConnection,
	}
}
