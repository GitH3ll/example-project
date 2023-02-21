package main

import (
	"fmt"
	"github.com/GitH3ll/example-project/internal/config"
	"github.com/GitH3ll/example-project/internal/repository"
	"github.com/GitH3ll/example-project/internal/server"
	"github.com/GitH3ll/example-project/internal/service"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()

	cfg := &config.Config{}

	err := cfg.Process()
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info(cfg.DB.Driver)

	db, err := sqlx.Connect(cfg.DB.Driver, fmt.Sprintf("user=%s dbname=%s sslmode=%s password=%s", cfg.DB.User,
		cfg.DB.Name, cfg.DB.SSLMode, cfg.DB.Password))
	if err != nil {
		logger.Fatal(err)
	}

	userRepo := repository.NewUserRepo(db, cfg.DB)

	err = userRepo.RunMigrations()
	if err != nil {
		logger.Warning(err)
	}

	controller := service.NewController(userRepo)

	srv := server.NewServer(":8000", logger, controller)
	srv.RegisterRoutes()

	srv.StartServer()
}
