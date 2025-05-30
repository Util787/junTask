package config

import (
	"os"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Env  string
	Port string
}

func InitServerConfig() *ServerConfig {
	err := godotenv.Load()
	if err != nil {
		panic("Failed to initialize server config. Make sure all required .env variables are set")
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
		Host:     os.Getenv("DBHOST"),
		Port:     os.Getenv("DBPORT"),
		Username: os.Getenv("DBUSERNAME"),
		Password: os.Getenv("DBPASSWORD"),
		DBName:   os.Getenv("DBNAME"),
		SSLMode:  os.Getenv("SSLMODE"),
	}
}
