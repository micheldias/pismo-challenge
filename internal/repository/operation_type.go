package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"pismo-challenge/internal/model"
)

type OperationType interface {
	GetById(context.Context, int) (model.OperationType, error)
}

func NewOperationType(database *sql.DB) OperationType {
	return &operationType{
		database: database,
	}
}

type operationType struct {
	database *sql.DB
}

func (a operationType) GetById(ctx context.Context, id int) (model.OperationType, error) {
	var result model.OperationType
	querySelect := `SELECT operation_type_id, description, direction  FROM operation_types where operation_type_id = $1`
	if err := a.database.QueryRowContext(ctx, querySelect, id).Scan(&result.Id, &result.Description, &result.Direction); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.OperationType{}, model.EntityNorFoundError{Reason: "operation type not found"}
		}

		return model.OperationType{}, fmt.Errorf("error to find operation type with id %d to the database: %w", id, err)
	}
	return result, nil
}

// --- INTERFACE GUARD ---
var _ OperationType = (*operationType)(nil)
