package services

type Config interface {
	PostgresDbName() string
	PostgresHost() string
	PostgresPassword() string
	PostgresUser() string
	PostgresPort() string
	HttpPort() string
	AccessSecret() string
	RefreshSecret() string
	AccessExpire() string
	RefreshExpire() string
	PublicPemPath() string
	PrivatePemPath() string
	NotificationGrpcPort() string
	NotificationGrpcHost() string
	ParkingGrpcPort() string
	ParkingGrpcHost() string
	AdminPassword() string
	AdminLogin() string
}
