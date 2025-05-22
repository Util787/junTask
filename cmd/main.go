package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Util787/junTask/config"
	"github.com/Util787/junTask/entities"
	"github.com/Util787/junTask/internal/database"
	"github.com/Util787/junTask/internal/handlers"
	"github.com/Util787/junTask/internal/repository"
	service "github.com/Util787/junTask/internal/services"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	servConfig, err := config.InitServerConfig()
	if err != nil {
		logrus.Fatal("Failed to initialize server config. Make sure all required .env variables are set.")
	}

	dbConfig := config.InitDbConfig()

	postgresDB, err := database.NewPostgresDB(*dbConfig)
	if err != nil {
		logrus.Fatal("Failed to connect to db")
	}
	repos := repository.NewRepository(postgresDB)
	services := service.NewService(repos)
	handlers := handlers.NewHandlers(services)

	srv := entities.Server{}
	go func() {
		err := srv.Run(servConfig.Port, handlers.InitRoutes())
		if err != nil {
			logrus.Fatal("Failed to run the server: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("Server shutting down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Println("Failed to shut down the server: ", err)
	}

	if err := postgresDB.Close(); err != nil {
		logrus.Println("error occured on db connection close: ", err)
	}
}
