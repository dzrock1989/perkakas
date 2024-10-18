package constant

// Static messages
const (
	MSG_FORBIDDEN_ACCESS     = "forbidden access"
	MSG_PERMANENTLY_REDIRECT = "permanently_redirect"
	MSG_NOT_FOUND            = "not found"
	MSG_SUCCESS              = "success"
)

// Env Keys
const (
	ENV            = "ENV"
	APP_PORT       = "APP_PORT"
	GRPC_PORT      = "GRPC_PORT"
	WEB_URL        = "WEB_URL"
	WEB_URL_YANMAS = "WEB_URL_YANMAS"

	IS_USE_ELK = "IS_USE_ELK"

	DB_HOST          = "DB_HOST"
	DB_PORT          = "DB_PORT"
	DB_USER          = "DB_USER"
	DB_PASS          = "DB_PASS"
	DB_NAME          = "DB_NAME"
	DB_MAX_OPEN_CONN = "DB_MAX_OPEN_CONN"
	DB_MAX_IDLE_CONN = "DB_MAX_IDLE_CONN"
	DB_MAX_LIFE_TIME = "DB_MAX_LIFE_TIME"

	// redis
	REDIS_ENABLED  = "REDIS_ENABLED"
	REDIS_HOST     = "REDIS_HOST"
	REDIS_PORT     = "REDIS_PORT"
	REDIS_USER     = "REDIS_USER"
	REDIS_PASS     = "REDIS_PASS"
	REDIS_DB       = "REDIS_DB"
	REDIS_AUTH_KEY = "REDIS_AUTH_KEY"

	REDIS_IS_CLUSTER       = "REDIS_IS_CLUSTER"
	REDIS_HOSTS            = "REDIS_HOSTS"
	REDIS_POOL_SIZE        = "REDIS_POOL_SIZE"
	REDIS_MAX_ACTIVE_CONNS = "REDIS_MAX_ACTIVE_CONNS"

	// auth
	JWT_SECRET_KEY = "JWT_SECRET_KEY"
	HMAC_DATE_KEY  = "HMAC_DATE_KEY"

	YANMAS_JWT_SECRET_KEY = "YANMAS_JWT_SECRET_KEY"
	YANMAS_HMAC_DATE_KEY  = "YANMAS_HMAC_DATE_KEY"

	VAULT_ENABLED      = "VAULT_ENABLED"
	VAULT_ADDRESS      = "VAULT_ADDRESS"
	VAULT_TOKEN        = "VAULT_TOKEN"
	VAULT_SERVICE_NAME = "SERVICE_NAME"
	VAULT_SECRET_NAME  = "SECRET_NAME"

	NATS_URL = "NATS_URL"

	ALLOWED_ORIGINS = "ALLOWED_ORIGINS"

	// Redpanda
	RP_HOST           = "RP_HOST"
	RP_PORT           = "RP_PORT"
	RP_REPLICA_FACTOR = "RP_REPLICA_FACTOR" // sesuai dengan jumlah redpanda instance pada cluster
	RP_GROUP          = "RP_GROUP"
	RP_LOG            = "RP_LOG" // true or false, default false
	// eg: topic1,topic2,topic3, default : foo
	RP_TOPICS = "RP_TOPICS"

	// SMTP
	SMTP_HOST   = "SMTP_HOST"
	SMTP_PORT   = "SMTP_PORT"
	SMTP_SENDER = "SMTP_SENDER"
	SMTP_USER   = "SMTP_USER"
	SMTP_PASS   = "SMTP_PASS"

	ES_URL           = "ES_URL"
	ES_API_KEY       = "ES_API_KEY"
	ES_INDEX_HISTORY = "ES_INDEX_HISTORY"
)
