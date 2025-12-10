package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Postgres struct {
	Port     string `env:"POSTGRES_PORT"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Database string `env:"NOTIFICATION_DB_NAME"`
	Host     string `env:"POSTGRES_HOST"`
}

type HttpServer struct {
	Port string `env:"NOTIFICATION_HTTP_PORT"`
}

type Smtp struct {
	Host string `env:"SMTP_HOST"`
	Port string `env:"SMTP_PORT"`
	User string `env:"SMTP_USER"`
	Pass string `env:"SMTP_PASS"`
	From string `env:"SMTP_FROM"`
}

type NotificationGrpc struct {
	Port string `env:"NOTIFICATION_GRPC_PORT"`
	Host string `env:"NOTIFICATION_GRPC_HOST"`
}

type Elastic struct {
	Host     string `env:"ELASTIC_HOST"`
	Port     string `env:"ELASTIC_PORT"`
	Username string `env:"ELASTIC_USERNAME"`
	Password string `env:"ELASTIC_PASSWORD"`
}

type AMQP struct {
	User     string `env:"RABBIT_USER"`
	Password string `env:"RABBIT_PASS"`
	Host     string `env:"RABBIT_HOST"`
	Port     string `env:"RABBIT_PORT"`
}

type Config struct {
	Postgres         Postgres
	HttpServer       HttpServer
	Smtp             Smtp
	NotificationGrpc NotificationGrpc
	Elastic          Elastic
	AMQP             AMQP
}

func Make() *Config {
	if err := godotenv.Load(".env"); err != nil {
		panic(".env file not found")
	}

	var config Config
	if err := cleanenv.ReadEnv(&config); err != nil {
		panic("cannot read environment")
	}

	return &config
}

func (c *Config) PostgresPort() string {
	return c.Postgres.Port
}

func (c *Config) PostgresUser() string {
	return c.Postgres.User
}

func (c *Config) PostgresPassword() string {
	return c.Postgres.Password
}

func (c *Config) PostgresDbName() string {
	return c.Postgres.Database
}

func (c *Config) PostgresHost() string {
	return c.Postgres.Host
}

func (c *Config) HttpPort() string {
	return c.HttpServer.Port
}

func (c *Config) SmtpPort() string {
	return c.Smtp.Port
}

func (c *Config) SmtpHost() string {
	return c.Smtp.Host
}

func (c *Config) SmtpUser() string {
	return c.Smtp.User
}

func (c *Config) SmtpPassword() string {
	return c.Smtp.Pass
}

func (c *Config) SmtpFrom() string {
	return c.Smtp.From
}

func (c *Config) NotificationGrpcPort() string {
	return c.NotificationGrpc.Port
}

func (c *Config) NotificationGrpcHost() string {
	return c.NotificationGrpc.Host
}

func (c *Config) ElasticHost() string {
	return c.Elastic.Host
}

func (c *Config) ElasticPort() string {
	return c.Elastic.Port
}

func (c *Config) ElasticUsername() string {
	return c.Elastic.Username
}

func (c *Config) ElasticPassword() string {
	return c.Elastic.Password
}

func (c *Config) AMQPUser() string {
	return c.AMQP.User
}

func (c *Config) AMQPPassword() string {
	return c.AMQP.Password
}

func (c *Config) AMQPHost() string {
	return c.AMQP.Host
}

func (c *Config) AMQPPort() string {
	return c.AMQP.Port
}
