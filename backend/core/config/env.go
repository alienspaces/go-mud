package config

const (
	// General
	AppServerEnv = "APP_SERVER_ENV"

	AppServerHost = "APP_SERVER_HOST"
	AppServerHome = "APP_SERVER_HOME"
	AppServerPort = "APP_SERVER_PORT"

	// Logger
	AppServerLogLevel  = "APP_SERVER_LOG_LEVEL"
	AppServerLogPretty = "APP_SERVER_LOG_PRETTY"

	// Database
	AppServerDBHost               = "APP_SERVER_DB_HOST"
	AppServerDBPort               = "APP_SERVER_DB_PORT"
	AppServerDBName               = "APP_SERVER_DB_NAME"
	AppServerDBUser               = "APP_SERVER_DB_USER"
	AppServerDBPassword           = "APP_SERVER_DB_PASSWORD"
	AppServerDBMaxOpenConnections = "APP_SERVER_DB_MAX_OPEN_CONNECTIONS"
	AppServerDBMaxIdleConnections = "APP_SERVER_DB_MAX_IDLE_CONNECTIONS"
	AppServerDBMaxIdleTimeMins    = "APP_SERVER_DB_MAX_IDLE_TIME_MINS"

	// SMTP
	AppServerSMTPHost = "APP_SERVER_SMTP_HOST"

	// For testing
	AppServerTxRollback = "APP_SERVER_API_SHOULD_DB_TX_ROLLBACK"
)

const (
	AppServerEnvDevelopValue    = "develop"
	AppServerEnvQAValue         = "qa"
	AppServerEnvStagingValue    = "staging"
	AppServerEnvProductionValue = "production"
)
