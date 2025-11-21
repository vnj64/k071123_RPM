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

type Middleware struct {
	PublicPemPath string `env:"JWT_PUBLIC_PEM_PATH"`
}

type Config struct {
	Postgres    Postgres
	HttpServer  HttpServer
	ParkingGrpc ParkingGrpc
	Middleware  Middleware
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
