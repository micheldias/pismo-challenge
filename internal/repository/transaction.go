package repository

import (
	"context"
	"database/sql"
	"fmt"
	"pismo-challenge/internal/model"
	"time"
)

type Transaction interface {
	Create(ctx context.Context, transaction model.Transaction) (model.Transaction, error)
}

func NewTransaction(database *sql.DB) Transaction {
	return &transaction{
		database: database,
	}
}

type transaction struct {
	database *sql.DB
}

func (t transaction) Create(ctx context.Context, transaction model.Transaction) (model.Transaction, error) {
	transaction.EventDate = time.Now().UTC()
	queryInsert := `INSERT INTO transactions (account_id, operation_type_id, amount, event_date) VALUES ($1, $2,$3,$4) RETURNING transaction_id`
	err := t.database.QueryRowContext(ctx, queryInsert, transaction.Account.Id, transaction.OperationTypes.Id, transaction.Amount, transaction.EventDate).Scan(&transaction.Id)
	if err != nil {
		return model.Transaction{}, fmt.Errorf("error saving transaction in account_id %d to the database: %w", transaction.Account.Id, err)
	}
	return transaction, nil
}

// --- INTERFACE GUARD ---
var _ Transaction = (*transaction)(nil)
