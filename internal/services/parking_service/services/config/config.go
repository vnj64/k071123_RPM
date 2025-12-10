package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Postgres struct {
	Port     string `env:"POSTGRES_PORT"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Database string `env:"PARKING_DB_NAME"`
	Host     string `env:"POSTGRES_HOST"`
}

type HttpServer struct {
	Port string `env:"PARKING_HTTP_PORT"`
}

type ParkingGrpc struct {
	Port string `env:"PARKING_GRPC_PORT"`
	Host string `env:"PARKING_GRPC_HOST"`
}

type OrderGrpc struct {
	Port string `env:"ORDER_GRPC_PORT"`
	Host string `env:"ORDER_GRPC_HOST"`
}

type UserGrpc struct {
	Port string `env:"USER_GRPC_PORT"`
	Host string `env:"USER_GRPC_HOST"`
}

type NotificationGrpc struct {
	Port string `env:"NOTIFICATION_GRPC_PORT"`
	Host string `env:"NOTIFICATION_GRPC_HOST"`
}

type Middleware struct {
	PublicPemPath string `env:"JWT_PUBLIC_PEM_PATH"`
}

type Elastic struct {
	Host     string `env:"ELASTIC_HOST"`
	Port     string `env:"ELASTIC_PORT"`
	Username string `env:"ELASTIC_USERNAME"`
	Password string `env:"ELASTIC_PASSWORD"`
}

type Config struct {
	Postgres         Postgres
	HttpServer       HttpServer
	ParkingGrpc      ParkingGrpc
	OrderGrpc        OrderGrpc
	UserGrpc         UserGrpc
	Middleware       Middleware
	Elastic          Elastic
	NotificationGrpc NotificationGrpc
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

func (c *Config) ParkingGrpcPort() string {
	return c.ParkingGrpc.Port
}

func (c *Config) ParkingGrpcHost() string {
	return c.ParkingGrpc.Host
}

func (c *Config) PublicPemPath() string {
	return c.Middleware.PublicPemPath
}

func (c *Config) OrderGrpcPort() string {
	return c.OrderGrpc.Port
}

func (c *Config) OrderGrpcHost() string {
	return c.OrderGrpc.Host
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

func (c *Config) UserGrpcHost() string {
	return c.UserGrpc.Host
}

func (c *Config) UserGrpcPort() string {
	return c.UserGrpc.Port
}

func (c *Config) NotificationGrpcPort() string {
	return c.NotificationGrpc.Port
}

func (c *Config) NotificationGrpcHost() string {
	return c.NotificationGrpc.Host
}
