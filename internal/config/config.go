package config

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Env               string        `env:"ENV" envDefault:"prod"`
	Port              string        `env:"HTTP_PORT" envDefault:"8000"`
	ReadHeaderTimeout time.Duration `env:"HTTP_READ_HEADER_TIMEOUT" envDefault:"5s"`
	WriteTimeout      time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"10s"`
	ReadTimeout       time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"10s"`
}

func InitServerConfig() *ServerConfig {
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load .env file. " + err.Error())
	}

	srvCfg := &ServerConfig{}

	if err := env.Parse(srvCfg); err != nil {
		panic("Failed to parse server config. " + err.Error())
	}

	if srvCfg.Env != "prod" && srvCfg.Env != "dev" && srvCfg.Env != "local" {
		panic("Invalid ENV variable, must be prod, dev or local")
	}

	return srvCfg
}

type DBConfig struct {
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     string `env:"DB_PORT" envDefault:"5432"`
	Username string `env:"DB_USERNAME" envDefault:"postgres"`
	Password string `env:"DB_PASSWORD"`
	DBName   string `env:"DB_NAME" envDefault:"postgres"`
	SSLMode  string `env:"DB_SSLMODE" envDefault:"disable"`
}

// didnt use godotenv.Load() here, init only after InitServerConfig()!
func InitDbConfig() *DBConfig {
	dbCfg := &DBConfig{}

	if err := env.Parse(dbCfg); err != nil {
		panic("Failed to parse db config. " + err.Error())
	}

	if dbCfg.Password == "" {
		panic("DB_PASSWORD is not set")
	}

	return dbCfg
}

type RedisConfig struct {
	Host     string `env:"REDIS_HOST" envDefault:"localhost"`
	Port     string `env:"REDIS_PORT" envDefault:"6379"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB" envDefault:"0"`
}

func InitRedisConfig() *RedisConfig {
	redisCfg := &RedisConfig{}

	if err := env.Parse(redisCfg); err != nil {
		panic("Failed to parse redis config. " + err.Error())
	}

	if redisCfg.Password == "" {
		panic("REDIS_PASSWORD is not set")
	}

	return redisCfg
}
