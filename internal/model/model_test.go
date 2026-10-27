package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTransaction_Validate(t *testing.T) {
	tests := []struct {
		name          string
		transaction   Transaction
		expectedError error
	}{
		{
			name: "should return nil error when transaction is valid",
			transaction: Transaction{
				Account:        Account{Id: 1, DocumentNumber: "123456"},
				Amount:         100,
				OperationTypes: OperationType{Id: 1},
			},
			expectedError: nil,
		},
		{
			name:        "should return the list with all errors",
			transaction: Transaction{Id: 1, Amount: 1},
			expectedError: MissingRequiredInformation{
				Errors: []error{
					MissingRequiredInformation{Reason: "invalid account id"},
					MissingRequiredInformation{Reason: "invalid operation type"},
				}},
		},
		{
			name:          "should return error when transaction contains a negative amount",
			transaction:   Transaction{Id: 1, Account: Account{Id: 1, DocumentNumber: "123456"}, Amount: -1, OperationTypes: OperationType{Id: 1}, EventDate: time.Now()},
			expectedError: MissingRequiredInformation{Errors: []error{MissingRequiredInformation{Reason: "amount must be greater than zero"}}},
		},
		{
			name:          "should return error when transaction contains a zero amount",
			transaction:   Transaction{Id: 1, Account: Account{Id: 1, DocumentNumber: "123456"}, Amount: 0, OperationTypes: OperationType{Id: 1}, EventDate: time.Now()},
			expectedError: MissingRequiredInformation{Errors: []error{MissingRequiredInformation{Reason: "amount must be greater than zero"}}},
		},
		{
			name:          "should return error when transaction contains invalid account",
			transaction:   Transaction{Id: 1, Account: Account{}, Amount: 1, OperationTypes: OperationType{Id: 1}, EventDate: time.Now()},
			expectedError: MissingRequiredInformation{Errors: []error{MissingRequiredInformation{Reason: "invalid account id"}}},
		},
		{
			name:          "should return error when transaction contains invalid operation type",
			transaction:   Transaction{Id: 1, Account: Account{Id: 1}, Amount: 1, OperationTypes: OperationType{Id: -1}, EventDate: time.Now()},
			expectedError: MissingRequiredInformation{Errors: []error{MissingRequiredInformation{Reason: "invalid operation type"}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.transaction.Validate()
			assert.Equal(t, tt.expectedError, err)

		})
	}
}

func TestAccount_Validate(t *testing.T) {
	tests := []struct {
		name          string
		account       Account
		expectedError error
	}{
		{
			name:          "should return nil error when account is valid",
			account:       Account{Id: 1, DocumentNumber: "123456"},
			expectedError: nil,
		},
		{
			name:          "should return the list with all errors",
			account:       Account{},
			expectedError: MissingRequiredInformation{Reason: "Account [document_number] cannot be empty"},
		},
		{
			name:          "should return error when document number is empty",
			account:       Account{Id: 1, DocumentNumber: ""},
			expectedError: MissingRequiredInformation{Reason: "Account [document_number] cannot be empty"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.account.Validate()
			assert.Equal(t, tt.expectedError, err)

		})
	}
}
