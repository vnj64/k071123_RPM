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

	ElasticPassword() string
	ElasticUsername() string
	ElasticPort() string
	ElasticHost() string
	AMQPPassword() string
	AMQPHost() string
	AMQPPort() string
	AMQPUser() string
}
