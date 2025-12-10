package logger

import (
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-extras/elogrus.v8"
	"net/http"
	"strings"
	"time"
)

type Logger struct {
	*logrus.Logger
}

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Index    string
	Service  string
}

func New(cfg Config) (*Logger, error) {
	if cfg.Host == "" || cfg.Port == "" || cfg.Index == "" {
		return nil, errors.New("invalid logger config: host/port/index required")
	}

	base := logrus.New()
	base.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	})
	base.SetLevel(logrus.InfoLevel)

	esURL := fmt.Sprintf("http://%s:%s", cfg.Host, cfg.Port)

	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{esURL},
		Username:  cfg.Username,
		Password:  cfg.Password,
		Transport: &http.Transport{},
	})
	if err != nil {
		return nil, err
	}

	// Проверка индекса
	exists, err := esClient.Indices.Exists([]string{cfg.Index})
	if err != nil {
		return nil, err
	}
	if exists.StatusCode == 404 {
		_, err := esClient.Indices.Create(
			cfg.Index,
			esClient.Indices.Create.WithBody(strings.NewReader(
				`{
				  "mappings": {
					"properties": {
					  "@timestamp": { "type": "date" },
					  "level": { "type": "keyword" },
					  "service": { "type": "keyword" },
					  "message": { "type": "text" }
					}
				  }
				}`,
			)),
		)
		if err != nil {
			return nil, err
		}
	}

	// Hook
	hook, err := elogrus.NewAsyncElasticHook(esClient, esURL, logrus.DebugLevel, cfg.Index)
	if err != nil {
		return nil, err
	}
	base.Hooks.Add(hook)

	return &Logger{Logger: base}, nil
}
