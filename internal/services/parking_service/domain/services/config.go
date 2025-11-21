package services

type Config interface {
	PostgresDbName() string
	PostgresHost() string
	PostgresPassword() string
	PostgresUser() string
	PostgresPort() string
	HttpPort() string
	ParkingGrpcPort() string
	ParkingGrpcHost() string
}
