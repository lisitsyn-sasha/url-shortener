package constants

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

const (
	StatusOk    = "OK"
	StatusError = "Error"
)

const (
	UrlSaveMew = "handlers.url.save.New"

	PostgresNew         = "storage.postgres.New"
	PostgresCreateTable = "storage.postgres.createTables"
	PostgresSaveUrl     = "storage.postgres.SaveURL"
	PostgresGetUrl      = "storage.postgres.GetURL"
)

const AliasLength = 6

const Charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
