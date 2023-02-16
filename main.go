package main

import (
	"github.com/GitH3ll/example-project/internal/repository"
	"github.com/GitH3ll/example-project/internal/server"
	"github.com/GitH3ll/example-project/internal/service"
	"github.com/sirupsen/logrus"
)

func main() {
	userRepo := repository.NewUserRepo("myNewFile.json")
	controller := service.NewController(userRepo)

	logger := logrus.New()

	srv := server.NewServer(":8000", logger, controller)
	srv.RegisterRoutes()

	srv.StartServer()
}
