package config

import (
	"os"

	"github.com/spf13/cast"
)

type Config struct {
	HTTP_PORT     string
	GRPC_PORT     string
	DB_HOST       string
	DB_PORT       int
	DB_USER       string
	DB_NAME       string
	DB_PASSWORD   string
	ACCESS_TOKEN  string
	REFRESH_TOKEN string
}

func coalasce(env string, defaultValue interface{}) interface{} {
	value, exists := os.LookupEnv(env)
	if !exists {
		return defaultValue
	}
	return value
}

func Load() Config {
	cfg := Config{}
	cfg.HTTP_PORT = cast.ToString(coalasce("HTTP_PORT","auth-service:8081"))
	cfg.GRPC_PORT = cast.ToString(coalasce("GRPC_PORT","auth-service:50053"))
	cfg.DB_HOST = cast.ToString(coalasce("DB_HOST","postgres-db"))
	cfg.DB_PORT = cast.ToInt(coalasce("DB_PORT",5432))
	cfg.DB_USER = cast.ToString(coalasce("DB_USER","postgres"))
	cfg.DB_NAME = cast.ToString(coalasce("DB_NAME","car_auth"))
	cfg.DB_PASSWORD = cast.ToString(coalasce("DB_PASSWORD","1918"))
	cfg.ACCESS_TOKEN = cast.ToString(coalasce("ACCESS_TOKEN","my_secret_key"))
	cfg.REFRESH_TOKEN = cast.ToString(coalasce("REFRESH_TOKEN","my_secret_key"))

	return cfg
}
