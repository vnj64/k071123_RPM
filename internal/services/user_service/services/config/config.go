package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Postgres struct {
	Port     string `env:"POSTGRES_PORT"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Database string `env:"USER_DB_NAME"`
	Host     string `env:"POSTGRES_HOST"`
}

type HttpServer struct {
	Port string `env:"USER_HTTP_PORT"`
}

type Jwt struct {
	AccessSecret   string `env:"JWT_ACCESS_SECRET"`
	RefreshSecret  string `env:"JWT_REFRESH_SECRET"`
	AccessExpire   string `env:"JWT_ACCESS_EXPIRE"`
	RefreshExpire  string `env:"JWT_REFRESH_EXPIRE"`
	PublicPemPath  string `env:"JWT_PUBLIC_PEM_PATH"`
	PrivatePemPath string `env:"JWT_PRIVATE_PEM_PATH"`
}

type NotificationGrpc struct {
	Port string `env:"NOTIFICATION_GRPC_PORT"`
	Host string `env:"NOTIFICATION_GRPC_HOST"`
}

type ParkingGrpc struct {
	Port string `env:"PARKING_GRPC_PORT"`
	Host string `env:"PARKING_GRPC_HOST"`
}

type UserGrpc struct {
	Port string `env:"USER_GRPC_PORT"`
	Host string `env:"USER_GRPC_HOST"`
}

type AdminData struct {
	Login    string `env:"ADMIN_LOGIN"`
	Password string `env:"ADMIN_PASSWORD"`
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
	Jwt              Jwt
	NotificationGrpc NotificationGrpc
	ParkingGrpc      ParkingGrpc
	UserGrpc         UserGrpc
	AdminData        AdminData
	Elastic          Elastic
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

func (c *Config) AccessSecret() string {
	return c.Jwt.AccessSecret
}

func (c *Config) RefreshSecret() string {
	return c.Jwt.RefreshSecret
}

func (c *Config) AccessExpire() string {
	return c.Jwt.AccessExpire
}

func (c *Config) RefreshExpire() string {
	return c.Jwt.RefreshExpire
}

func (c *Config) PublicPemPath() string {
	return c.Jwt.PublicPemPath
}

func (c *Config) PrivatePemPath() string {
	return c.Jwt.PrivatePemPath
}

func (c *Config) NotificationGrpcPort() string {
	return c.NotificationGrpc.Port
}

func (c *Config) NotificationGrpcHost() string {
	return c.NotificationGrpc.Host
}

func (c *Config) ParkingGrpcPort() string {
	return c.ParkingGrpc.Port
}

func (c *Config) ParkingGrpcHost() string {
	return c.ParkingGrpc.Host
}

func (c *Config) AdminLogin() string {
	return c.AdminData.Login
}

func (c *Config) AdminPassword() string {
	return c.AdminData.Password
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

func (c *Config) UserGrpcPort() string {
	return c.UserGrpc.Port
}

func (c *Config) UserGrpcHost() string {
	return c.UserGrpc.Host
}
