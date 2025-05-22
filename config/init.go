package config

import (
	"os"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Port string
}

func InitServerConfig() (*ServerConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &ServerConfig{
		Port: os.Getenv("SERVERPORT"),
	}, nil
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
