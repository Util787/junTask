package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Env  string
	Port string
}

func InitServerConfig() *ServerConfig {
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load .env file. Make sure all required .env variables are set")
	}

	env := os.Getenv("ENV")
	if env != "local" && env != "prod" && env != "dev" {
		panic("Invalid environment value: must be 'prod', 'dev', or 'local'")
	}

	return &ServerConfig{
		Env:  env,
		Port: os.Getenv("SERVERPORT"),
	}
}

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

// didnt use godotenv.Load() here, init only after InitServerConfig()!
func InitDbConfig() *DBConfig {
	return &DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("SSLMODE"),
	}
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func InitRedisConfig() *RedisConfig {
	dbNum, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		panic("Invalid REDIS_DB value: must be digit")
	}
	return &RedisConfig{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       dbNum,
	}
}
