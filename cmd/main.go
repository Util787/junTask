package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/Util787/user-manager-api/config"
	"github.com/Util787/user-manager-api/entities"
	"github.com/Util787/user-manager-api/internal/handlers"
	"github.com/Util787/user-manager-api/internal/logger/handlers/slogpretty"
	"github.com/Util787/user-manager-api/internal/logger/sl"
	"github.com/Util787/user-manager-api/internal/repository"
	service "github.com/Util787/user-manager-api/internal/services"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// @title           User manager api
// @version         1.0
// @description     Rest api for managing users crud operations

// @host      localhost:8000
// @BasePath  /api

func main() {
	servConfig := config.InitServerConfig()

	log := setupLogger(servConfig.Env)

	//postgres init
	dbConfig := config.InitDbConfig()
	postgresDB, err := repository.NewPostgresDB(*dbConfig)
	if err != nil {
		log.Error("Failed to connect to db", sl.Err(err))
		return
	}
	log.Info("Connected to db successfully")

	//redis init
	redisConfig := config.InitRedisConfig()
	redis, err := repository.NewRedisClient(*redisConfig)
	if err != nil {
		log.Error("Failed to connect to redis", sl.Err(err))
		return
	}
	log.Info("Connected to redis successfully")

	//layers
	repos := repository.NewRepository(postgresDB, redis)
	services := service.NewService(repos)
	handlers := handlers.NewHandlers(services, log)

	//server start
	srv := entities.Server{}
	go func() {
		err := srv.CreateAndRun(servConfig.Port, handlers.InitRoutes(servConfig.Env))
		if err != nil {
			log.Error("Server was interrupted", sl.Err(err))
		}
	}()

	//graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Info("Shutting down the server")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Error("Failed to shut down the server", sl.Err(err))
	}

	if err := redis.Close(); err != nil {
		log.Error("Error occurred during redis connection closing", sl.Err(err))
	}

	if err := postgresDB.Close(); err != nil {
		log.Error("Error occured during db connection closing", sl.Err(err))
	}
	log.Info("Gracefully stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
