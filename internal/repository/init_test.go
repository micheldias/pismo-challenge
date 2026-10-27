package repository

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"pismo-challenge/internal/infra/database"
	"strconv"
	"testing"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

const (
	databaseName = "testdb"
	username     = "postgres"
	password     = "password"
)

var transactionRepository Transaction
var accountRepository Account
var operationRepository OperationType

func TestMain(m *testing.M) {
	code := runTests(m)
	os.Exit(code)
}

func runTests(m *testing.M) int {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(databaseName),
		postgres.WithUsername(username),
		postgres.WithPassword(password),
		postgres.WithInitScripts(filepath.Join("..", "..", "mocks", "populate.sql")),
		postgres.BasicWaitStrategies(),
	)

	if err != nil {
		log.Fatalf("Failed to start the container: %s", err)
	}
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Printf("Failed to terminate container: %s", err)
		}
	}()

	host, err := pgContainer.Host(ctx)
	if err != nil {
		log.Fatalf("Failed to get the host container: %s", err)
	}
	mappedPort, err := pgContainer.MappedPort(ctx, "5432/tcp")
	if err != nil {
		log.Fatalf("Failed to get the host container: %s", err)
	}

	portInt, err := strconv.Atoi(mappedPort.Port())
	if err != nil {
		log.Fatalf("Failed to get the port container: %s", err)
	}

	databaseConnection, err := database.NewPostgresDatabaseConnection(host, portInt, username, password, databaseName)
	if err != nil {
		log.Fatalf("Failed to get the database connection: %s", err)

	}
	defer databaseConnection.Close()
	transactionRepository = NewTransaction(databaseConnection)
	accountRepository = NewAccount(databaseConnection)
	operationRepository = NewOperationType(databaseConnection)
	return m.Run()
}
