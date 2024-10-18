package configs

import (
	"os"
	"strings"

	"github.com/tigapilarmandiri/perkakas/common/constant"
	"github.com/tigapilarmandiri/perkakas/common/util"

	"github.com/joho/godotenv"
	"github.com/tigapilarmandiri/perkakas"
)

type Configs struct {
	Env          string `json:"env"`
	IsTesting    bool   `json:"testing"`
	AppPort      string `json:"app_port"`
	GRPCPort     string `json:"grpc_port"`
	WebURL       string `json:"web_url"`
	WebURLYanmas string `json:"web_url_yanmas"`
	IsUseELK     bool   `json:"is_use_elk"`

	// Database connection info
	DB ConnInfo `json:"db"`

	// Redis connection info
	Redis Redis `json:"redis"`

	// JWT
	JWT `json:"jwt"`

	// NATS
	NatsURL        string `json:"nats_url"`
	AllowedOrigins string `json:"allowed_origins"`

	// Redpanda
	Redpanda `json:"redpanda"`

	// Redpanda
	SMTP `json:"smtp"`

	Elastic Elastic `json:"elastic"`
}

type Elastic struct {
	Url          string `json:"url"`
	ApiKey       string `json:"api_key"`
	IndexHistory string `json:"index_history"`
}

type ConnInfo struct {
	Host string `json:"host"`
	Port string `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
	// Eg: Database name
	Name string `json:"name"`

	MaxOpenConn int `json:"max_open_conn"`
	MaxIdleConn int `json:"max_idle_conn"`
	MaxLifeTime int `json:"max_life_time"` // will convert to minutes
}

type Redis struct {
	// Define if redis is enabled
	Enabled bool `json:"enabled"`

	Host         string `json:"host"`
	Port         string `json:"port"`
	User         string `json:"user"`
	Pass         string `json:"pass"`
	RedisAuthKey string `json:"redis_auth_key"`
	DB           int    `json:"db"`

	IsCluster      bool   `json:"is_cluster"`
	Hosts          string `json:"hosts"`
	PoolSize       int    `json:"pool_size"`
	MaxActiveConns int    `json:"max_active_conns"`
}

type JWT struct {
	SecretKey string `json:"secret_key"`
	DateKey   string `json:"date_key"`

	YanmasSecretKey string `json:"yanmas_secret_key"`
	YanmasDateKey   string `json:"yanmas_date_key"`
}

type ConfigOpts struct {
	// Contain path to .env file
	EnvFile string
}

type Redpanda struct {
	Host          string `json:"host"`
	Port          string `json:"port"`
	ReplicaFactor int    `json:"replica_factor"`
	Group         string `json:"group"`
	Topics        string `json:"topics"`
	Log           bool   `json:"log"`
}

type SMTP struct {
	Host   string `json:"host"`
	Port   int    `json:"port"`
	Sender string `json:"sender"`
	User   string `json:"user"`
	Pass   string `json:"pass"`
}

var Config Configs

func (conf Configs) IsDevelopment() bool {
	return strings.ToLower(conf.Env) == "development"
}

func (conf Configs) IsStaging() bool {
	return strings.ToLower(conf.Env) == "staging"
}

func (conf Configs) IsProduction() bool {
	return strings.ToLower(conf.Env) == "production"
}

var configOpts ConfigOpts

// Load configs from env or vault
func LoadConfigs() (confs Configs) {
	defer util.InitElk()
	if configOpts.EnvFile == "" {
		configOpts.EnvFile = ".env"
	}
	if err := godotenv.Load(configOpts.EnvFile); err != nil {
		util.Log.Info().Msg("No .env file found, will take host variable instead")
	}

	isConfigFromVault := perkakas.DefaultValueBoolFromString(false, os.Getenv(constant.VAULT_ENABLED))
	if isConfigFromVault {
		vaultConfig := Vault{
			Address:     os.Getenv(constant.VAULT_ADDRESS),
			Token:       os.Getenv(constant.VAULT_TOKEN),
			ServiceName: os.Getenv(constant.VAULT_SERVICE_NAME),
			SecretName:  os.Getenv(constant.VAULT_SECRET_NAME),
		}

		Config = vaultConfig.load()
		return
	}

	Config = Configs{
		Env:          perkakas.DefaultValueString("development", os.Getenv(constant.ENV)),
		AppPort:      perkakas.DefaultValueString("8080", os.Getenv(constant.APP_PORT)),
		GRPCPort:     perkakas.DefaultValueString("8081", os.Getenv(constant.GRPC_PORT)),
		WebURL:       perkakas.DefaultValueString("https://web.thefalcon.site", os.Getenv(constant.WEB_URL)),
		WebURLYanmas: perkakas.DefaultValueString("https://web.thefalcon.site", os.Getenv(constant.WEB_URL_YANMAS)),
		IsUseELK:     perkakas.DefaultValueBoolFromString(false, os.Getenv(constant.IS_USE_ELK)),
		DB: ConnInfo{
			Host:        perkakas.DefaultValueString("localhost", os.Getenv(constant.DB_HOST)),
			Port:        perkakas.DefaultValueString("5432", os.Getenv(constant.DB_PORT)),
			User:        perkakas.DefaultValueString("postgres", os.Getenv(constant.DB_USER)),
			Pass:        perkakas.DefaultValueString("secret", os.Getenv(constant.DB_PASS)),
			Name:        perkakas.DefaultValueString("postgres", os.Getenv(constant.DB_NAME)),
			MaxOpenConn: perkakas.DefaultValueIntFromString(100, os.Getenv(constant.DB_MAX_OPEN_CONN)),
			MaxIdleConn: perkakas.DefaultValueIntFromString(5, os.Getenv(constant.DB_MAX_IDLE_CONN)),
			MaxLifeTime: perkakas.DefaultValueIntFromString(15, os.Getenv(constant.DB_MAX_LIFE_TIME)),
		},
		Redis: Redis{
			Enabled:      perkakas.DefaultValueBoolFromString(true, os.Getenv(constant.REDIS_ENABLED)),
			Host:         perkakas.DefaultValueString("localhost", os.Getenv(constant.REDIS_HOST)),
			Port:         perkakas.DefaultValueString("6379", os.Getenv(constant.REDIS_PORT)),
			User:         perkakas.DefaultValueString("", os.Getenv(constant.REDIS_USER)),
			Pass:         perkakas.DefaultValueString("", os.Getenv(constant.REDIS_PASS)),
			DB:           perkakas.DefaultValueIntFromString(0, os.Getenv(constant.REDIS_DB)),
			RedisAuthKey: perkakas.DefaultValueString("all_roles", os.Getenv(constant.REDIS_AUTH_KEY)),

			IsCluster:      perkakas.DefaultValueBoolFromString(false, os.Getenv(constant.REDIS_IS_CLUSTER)),
			Hosts:          perkakas.DefaultValueString("", os.Getenv(constant.REDIS_HOSTS)),
			PoolSize:       perkakas.DefaultValueIntFromString(10, os.Getenv(constant.REDIS_POOL_SIZE)),
			MaxActiveConns: perkakas.DefaultValueIntFromString(10, os.Getenv(constant.REDIS_MAX_ACTIVE_CONNS)),
		},
		JWT: JWT{
			SecretKey:       perkakas.DefaultValueString("secretJwtKey", os.Getenv(constant.JWT_SECRET_KEY)),
			DateKey:         perkakas.DefaultValueString("secretDateKey", os.Getenv(constant.HMAC_DATE_KEY)),
			YanmasSecretKey: perkakas.DefaultValueString("secretJwtKey", os.Getenv(constant.YANMAS_JWT_SECRET_KEY)),
			YanmasDateKey:   perkakas.DefaultValueString("secretDateKey", os.Getenv(constant.YANMAS_HMAC_DATE_KEY)),
		},
		NatsURL:        perkakas.DefaultValueString("localhost:4222", os.Getenv(constant.NATS_URL)),
		AllowedOrigins: perkakas.DefaultValueString("*", os.Getenv(constant.ALLOWED_ORIGINS)),
		Redpanda: Redpanda{
			Host:          perkakas.DefaultValueString("localhost", os.Getenv(constant.RP_HOST)),
			Port:          perkakas.DefaultValueString("9092", os.Getenv(constant.RP_PORT)),
			ReplicaFactor: perkakas.DefaultValueIntFromString(1, os.Getenv(constant.RP_REPLICA_FACTOR)),
			Group:         perkakas.DefaultValueString("service-a", os.Getenv(constant.RP_GROUP)),
			Log:           perkakas.DefaultValueBoolFromString(false, os.Getenv(constant.RP_LOG)),
			Topics:        perkakas.DefaultValueString("foo", os.Getenv(constant.RP_TOPICS)),
		},
		SMTP: SMTP{
			Host:   perkakas.DefaultValueString("", os.Getenv(constant.SMTP_HOST)),
			Port:   perkakas.DefaultValueIntFromString(465, os.Getenv(constant.SMTP_PORT)),
			Sender: perkakas.DefaultValueString("", os.Getenv(constant.SMTP_SENDER)),
			User:   perkakas.DefaultValueString("", os.Getenv(constant.SMTP_USER)),
			Pass:   perkakas.DefaultValueString("", os.Getenv(constant.SMTP_PASS)),
		},

		Elastic: Elastic{
			Url:          os.Getenv(constant.ES_URL),
			ApiKey:       os.Getenv(constant.ES_API_KEY),
			IndexHistory: os.Getenv(constant.ES_INDEX_HISTORY),
		},
	}

	return
}

// Load configs from env or vault
func LoadConfigsWithOption(opts ConfigOpts) (confs Configs) {
	configOpts.EnvFile = opts.EnvFile

	return LoadConfigs()
}
