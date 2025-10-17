package services

type Config interface {
	PostgresDbName() string
	PostgresHost() string
	PostgresPassword() string
	PostgresUser() string
	PostgresPort() string
	HttpPort() string

	SmtpHost() string
	SmtpPort() string
	SmtpUser() string
	SmtpPassword() string
	SmtpFrom() string
}
