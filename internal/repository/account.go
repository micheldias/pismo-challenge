package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"pismo-challenge/internal/model"
)

type Account interface {
	Create(context.Context, model.Account) (model.Account, error)
	GetById(context.Context, int) (model.Account, error)
}

func NewAccount(database *sql.DB) Account {
	return &account{
		database: database,
	}
}

type account struct {
	database *sql.DB
}

func (a account) Create(ctx context.Context, account model.Account) (model.Account, error) {
	queryInsert := `INSERT INTO accounts (document_number) VALUES ($1) RETURNING account_id`
	err := a.database.QueryRowContext(ctx, queryInsert, account.DocumentNumber).Scan(&account.Id)
	if err != nil {
		return model.Account{}, fmt.Errorf("error saving account to the database: %w", err)
	}
	return account, nil
}

func (a account) GetById(ctx context.Context, id int) (model.Account, error) {
	var result model.Account
	querySelect := `SELECT account_id, document_number FROM accounts where account_id = $1`
	if err := a.database.QueryRowContext(ctx, querySelect, id).Scan(&result.Id, &result.DocumentNumber); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Account{}, model.EntityNorFoundError{Reason: "account not found"}
		}

		return model.Account{}, fmt.Errorf("error to find account with id %d to the database: %w", id, err)
	}
	return result, nil
}

// --- INTERFACE GUARD ---
var _ Account = (*account)(nil)
