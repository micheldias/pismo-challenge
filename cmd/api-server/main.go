package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	_ "pismo-challenge/docs"
	"pismo-challenge/internal/controller"
	"pismo-challenge/internal/infra/database"
	"pismo-challenge/internal/infra/logging"
	"pismo-challenge/internal/infra/server"
	servermodel "pismo-challenge/internal/infra/server/model"
	"pismo-challenge/internal/infra/settings"
	"pismo-challenge/internal/repository"
	"pismo-challenge/internal/service"

	"github.com/kelseyhightower/envconfig"
)

// @title			Pismo Challenge API
// @version		1.0
// @description	This is a  REST API with Chi and Swaggo
// @contact.name	Michel Dias
// @contact.email	michelc.dias@gmail.com
// @host			localhost:8080
// @BasePath		/
func main() {
	logging.ConfigureStructuredLog(os.Stdout)
	if err := envconfig.Process("app", &settings.Settings); err != nil {
		slog.Error(fmt.Sprintf("can't start env vars. error: %s", err.Error()))
		os.Exit(1)
	}
	databaseConnection, err := database.NewPostgresDatabaseConnection(settings.Settings.DatabaseConnection.Host, settings.Settings.DatabaseConnection.Port,
		settings.Settings.DatabaseConnection.User, settings.Settings.DatabaseConnection.Password, settings.Settings.DatabaseConnection.Dbname)
	if err != nil {
		slog.Error(fmt.Sprintf("It was not possible to establish a connection to the database: %s", err.Error()))
		os.Exit(1)
	}

	accountRepository := repository.NewAccount(databaseConnection)
	accountService := service.NewAccountService(accountRepository)
	accountController := controller.NewAccount(accountService)

	transactionService := service.NewTransactionService(repository.NewTransaction(databaseConnection), repository.NewOperationType(databaseConnection), accountRepository)
	transactionController := controller.NewTransaction(transactionService)

	serverApi := server.CreateWebServer(controller.GetRoutes(accountController, transactionController), []servermodel.ShutdownTask{
		{
			Name: "database",
			Task: func(ctx context.Context) error {
				slog.Info("shutting down database...")
				return databaseConnection.Close()
			},
		},
	})
	serverApi.Start()
}
