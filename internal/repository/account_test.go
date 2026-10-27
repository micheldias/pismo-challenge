package repository

import (
	"context"
	"pismo-challenge/internal/model"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccount_Create(t *testing.T) {
	tests := []struct {
		name           string
		input          model.Account
		expectedAssert func(t *testing.T, result model.Account, err error)
	}{
		{
			name:  "create account successfully",
			input: model.Account{DocumentNumber: "3456789"},
			expectedAssert: func(t *testing.T, result model.Account, err error) {
				assert.Nil(t, err)
				assert.NotZero(t, result.Id)
				assert.Equal(t, "3456789", result.DocumentNumber)
			},
		},
		{
			name:  "error when saving account with a document that exceeds the maximum limit",
			input: model.Account{DocumentNumber: strings.Repeat("9", 101)},
			expectedAssert: func(t *testing.T, result model.Account, err error) {
				assert.ErrorContains(t, err, "error saving account to the database")
				assert.Empty(t, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := accountRepository.Create(context.Background(), tt.input)
			tt.expectedAssert(t, result, err)
		})
	}
}

func TestAccount_GetById(t *testing.T) {
	tests := []struct {
		name           string
		input          int
		expectedAssert func(t *testing.T, result model.Account, err error)
	}{
		{
			name:  "get account by id successfully",
			input: 1,
			expectedAssert: func(t *testing.T, result model.Account, err error) {
				assert.Nil(t, err)
				assert.Equal(t, 1, result.Id)
				assert.Equal(t, "123456", result.DocumentNumber)
			},
		},
		{
			name:  "error when account does not exist",
			input: 99999,
			expectedAssert: func(t *testing.T, result model.Account, err error) {
				assert.ErrorAs(t, err, &model.EntityNorFoundError{})
				assert.ErrorContains(t, err, "account not found")
				assert.Empty(t, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := accountRepository.GetById(context.Background(), tt.input)
			tt.expectedAssert(t, result, err)
		})
	}
}
