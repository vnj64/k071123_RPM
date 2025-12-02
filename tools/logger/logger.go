package logger

import (
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-extras/elogrus.v8"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Index    string
}

var (
	instance *logrus.Logger
	once     sync.Once
)

func New(cfg Config) (*logrus.Logger, error) {
	var initErr error

	once.Do(func() {
		l := logrus.New()
		l.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
		l.SetLevel(logrus.DebugLevel)

		if cfg.Host == "" || cfg.Port == "" || cfg.Index == "" {
			initErr = errors.New("invalid logger config: host/port/index required")
			return
		}

		esURL := fmt.Sprintf("http://%s:%s", cfg.Host, cfg.Port)

		esClient, err := elasticsearch.NewClient(elasticsearch.Config{
			Addresses: []string{esURL},
			Username:  cfg.Username,
			Password:  cfg.Password,
			Transport: &http.Transport{Proxy: nil},
		})
		if err != nil {
			initErr = err
			return
		}

		exists, err := esClient.Indices.Exists([]string{cfg.Index})
		if err != nil {
			initErr = err
			return
		}

		if exists.StatusCode == 404 {
			_, err := esClient.Indices.Create(
				cfg.Index,
				esClient.Indices.Create.WithBody(strings.NewReader(
					`{
					  "mappings": {
						"properties": {
						  "@timestamp": { "type": "date" },
						  "message":    { "type": "text" }
						}
					  }
					}`,
				)),
			)

			if err != nil {
				initErr = err
				return
			}
		}

		hook, err := elogrus.NewAsyncElasticHook(esClient, esURL, logrus.DebugLevel, cfg.Index)
		if err != nil {
			initErr = err
			return
		}

		l.Hooks.Add(hook)

		instance = l
	})

	if initErr != nil {
		return nil, initErr
	}

	return instance, nil
}

func L() *logrus.Logger {
	if instance == nil {
		panic("logger not initialized: call logger.New() first")
	}
	return instance
}
