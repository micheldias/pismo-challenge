package repository

import (
	"context"
	"pismo-challenge/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionOperation_GetById(t *testing.T) {
	tests := []struct {
		name           string
		input          int
		expectedAssert func(t *testing.T, result model.OperationType, err error)
	}{
		{
			name:  "get operation type by id successfully",
			input: 1,
			expectedAssert: func(t *testing.T, result model.OperationType, err error) {
				assert.Nil(t, err)
				assert.Equal(t, 1, result.Id)
				assert.Equal(t, "Normal Purchase", result.Description)
				assert.Equal(t, "DEBIT", result.Direction)
			},
		},
		{
			name:  "error when operation type does not exist",
			input: 99999,
			expectedAssert: func(t *testing.T, result model.OperationType, err error) {
				assert.ErrorAs(t, err, &model.EntityNorFoundError{})
				assert.ErrorContains(t, err, "operation type not found")
				assert.Empty(t, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := operationRepository.GetById(context.Background(), tt.input)
			tt.expectedAssert(t, result, err)
		})
	}
}
