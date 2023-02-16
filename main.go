package main

import (
	"github.com/GitH3ll/example-project/internal/repository"
	"github.com/GitH3ll/example-project/internal/server"
	"github.com/GitH3ll/example-project/internal/service"
	"github.com/sirupsen/logrus"
)

func main() {
	fileName := "myfile.json"

	UserRepo := repository.NewUserRepository(fileName)

	controller := service.NewUserService(UserRepo)

	logger := logrus.New()

	srv := server.NewServer(logger, controller)
	srv.RegisterRoutes()

	srv.StartServer()
}
