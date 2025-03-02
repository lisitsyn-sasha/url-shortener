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
	UrlSaveNew   = "handlers.url.save.New"
	UrlDeleteNew = "handlers.url.delete.New"

	PostgresNew         = "storage.postgres.New"
	PostgresCreateTable = "storage.postgres.createTables"
	PostgresSaveUrl     = "storage.postgres.SaveURL"
	PostgresGetUrl      = "storage.postgres.GetURL"
	PostgresDeleteUrl   = "storage.postgres.DeleteURL"
)

const AliasLength = 6

const Charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
