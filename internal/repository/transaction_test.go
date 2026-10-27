package repository

import (
	"context"
	"pismo-challenge/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransaction_Create(t *testing.T) {
	tests := []struct {
		name           string
		input          model.Transaction
		expectedAssert func(t *testing.T, result model.Transaction, err error)
	}{
		{
			name: "create transaction successfully",
			input: model.Transaction{
				Account:        model.Account{Id: 1},
				OperationTypes: model.OperationType{Id: 1},
				Amount:         -12345,
			},
			expectedAssert: func(t *testing.T, result model.Transaction, err error) {
				assert.Nil(t, err)
				assert.NotZero(t, result.Id)
				assert.Equal(t, 1, result.OperationTypes.Id)
				assert.Equal(t, -12345, result.Amount)
				assert.NotZero(t, result.EventDate)
			},
		},
		{
			name: "error when saving transaction with a non-existent account",
			input: model.Transaction{
				Account:        model.Account{Id: 999},
				OperationTypes: model.OperationType{Id: 1},
				Amount:         -50,
			},
			expectedAssert: func(t *testing.T, result model.Transaction, err error) {
				assert.ErrorContains(t, err, "error saving transaction in account_id 999 to the database")
				assert.Empty(t, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := transactionRepository.Create(context.Background(), tt.input)
			tt.expectedAssert(t, result, err)
		})
	}
}
