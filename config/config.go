package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Port string
}

func InitConfig() *ServerConfig {
	err := godotenv.Load()
	if err != nil {
		log.Println("Failed to load .env file: ", err)
	}

	return &ServerConfig{
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

// didnt use godotenv.Load() here, init only after ServerConfig!
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
