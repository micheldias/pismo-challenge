package service

import (
	"context"
	"errors"
	"pismo-challenge/internal/model"
	"pismo-challenge/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransactionService_Create(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name                                string
		transaction                         model.Transaction
		transactionRepositoryMockReturned   func() (model.Transaction, error)
		operationTypeRepositoryMockReturned func() (model.OperationType, error)
		accountTypeRepositoryMockReturned   func() (model.Account, error)
		expectedResult                      func() (model.Transaction, error)
		expectedTransactionRepositoryCall   int
		expectedOperationTypeRepositoryCall int
		expectedAccountTypeRepositoryCall   int
	}{
		{
			name:                              "should forward the error received by the transaction repository",
			transaction:                       model.Transaction{Account: model.Account{Id: 1}, Amount: 1, OperationTypes: model.OperationType{Id: 4}},
			transactionRepositoryMockReturned: func() (model.Transaction, error) { return model.Transaction{}, errors.New("some kind of error") },
			operationTypeRepositoryMockReturned: func() (model.OperationType, error) {
				return model.OperationType{Id: 4, Direction: model.DirectionCredit}, nil
			},
			accountTypeRepositoryMockReturned:   func() (model.Account, error) { return model.Account{}, nil },
			expectedResult:                      func() (model.Transaction, error) { return model.Transaction{}, errors.New("some kind of error") },
			expectedTransactionRepositoryCall:   1,
			expectedOperationTypeRepositoryCall: 1,
			expectedAccountTypeRepositoryCall:   1,
		},
		{
			name:                                "should forward the error received by the operation type repository",
			transaction:                         model.Transaction{Account: model.Account{Id: 1}, Amount: 1, OperationTypes: model.OperationType{Id: 4}},
			transactionRepositoryMockReturned:   func() (model.Transaction, error) { return model.Transaction{}, nil },
			operationTypeRepositoryMockReturned: func() (model.OperationType, error) { return model.OperationType{}, errors.New("some kind of error") },
			accountTypeRepositoryMockReturned:   func() (model.Account, error) { return model.Account{}, nil },
			expectedResult:                      func() (model.Transaction, error) { return model.Transaction{}, errors.New("some kind of error") },
			expectedTransactionRepositoryCall:   0,
			expectedOperationTypeRepositoryCall: 1,
			expectedAccountTypeRepositoryCall:   0,
		},
		{
			name:                              "should forward the error received by the account repository",
			transaction:                       model.Transaction{Account: model.Account{Id: 1}, Amount: 1, OperationTypes: model.OperationType{Id: 4}},
			transactionRepositoryMockReturned: func() (model.Transaction, error) { return model.Transaction{}, nil },
			operationTypeRepositoryMockReturned: func() (model.OperationType, error) {
				return model.OperationType{Id: 4, Direction: model.DirectionCredit}, nil
			},
			accountTypeRepositoryMockReturned:   func() (model.Account, error) { return model.Account{}, errors.New("some kind of error") },
			expectedResult:                      func() (model.Transaction, error) { return model.Transaction{}, errors.New("some kind of error") },
			expectedTransactionRepositoryCall:   0,
			expectedOperationTypeRepositoryCall: 1,
			expectedAccountTypeRepositoryCall:   1,
		},
		{
			name:                              "should return the error without call the repository if an invalid info was received",
			transaction:                       model.Transaction{Account: model.Account{Id: 1}, OperationTypes: model.OperationType{Id: 4}},
			transactionRepositoryMockReturned: func() (model.Transaction, error) { return model.Transaction{}, nil },
			operationTypeRepositoryMockReturned: func() (model.OperationType, error) {
				return model.OperationType{Id: 4, Direction: model.DirectionCredit}, nil
			},
			accountTypeRepositoryMockReturned: func() (model.Account, error) { return model.Account{}, nil },
			expectedResult: func() (model.Transaction, error) {
				return model.Transaction{}, model.MissingRequiredInformation{Errors: []error{model.MissingRequiredInformation{Reason: "amount must be greater than zero"}}}
			},
			expectedTransactionRepositoryCall:   0,
			expectedOperationTypeRepositoryCall: 1,
			expectedAccountTypeRepositoryCall:   1,
		},
		{
			name:        "should create a transaction successfully",
			transaction: model.Transaction{Account: model.Account{Id: 1}, Amount: 1, OperationTypes: model.OperationType{Id: 4}},
			transactionRepositoryMockReturned: func() (model.Transaction, error) {
				return model.Transaction{Id: 1, Account: model.Account{Id: 1}, Amount: 1, OperationTypes: model.OperationType{Id: 4}}, nil
			},
			operationTypeRepositoryMockReturned: func() (model.OperationType, error) {
				return model.OperationType{Id: 4, Direction: model.DirectionCredit}, nil
			},
			accountTypeRepositoryMockReturned: func() (model.Account, error) { return model.Account{}, nil },
			expectedResult: func() (model.Transaction, error) {
				return model.Transaction{Id: 1, Account: model.Account{Id: 1}, Amount: 1, OperationTypes: model.OperationType{Id: 4}}, nil
			},
			expectedTransactionRepositoryCall:   1,
			expectedOperationTypeRepositoryCall: 1,
			expectedAccountTypeRepositoryCall:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transactionRepository := new(TransactionMock)
			operationTypeRepository := new(OperationTypeMock)
			accountRepository := new(AccountRepositoryMock)

			transactionRepository.On("Create", ctx, tt.transaction).Return(tt.transactionRepositoryMockReturned())
			operationTypeRepository.On("GetById", ctx, tt.transaction.OperationTypes.Id).Return(tt.operationTypeRepositoryMockReturned())
			accountRepository.On("GetById", ctx, tt.transaction.Account.Id).Return(tt.accountTypeRepositoryMockReturned())
			transactionCreated, err := NewTransactionService(transactionRepository, operationTypeRepository, accountRepository).Create(ctx, tt.transaction)
			expectedTransaction, expectedError := tt.expectedResult()
			assert.Equal(t, expectedTransaction, transactionCreated)
			assert.Equal(t, expectedError, err)
			transactionRepository.AssertNumberOfCalls(t, "Create", tt.expectedTransactionRepositoryCall)
			operationTypeRepository.AssertNumberOfCalls(t, "GetById", tt.expectedOperationTypeRepositoryCall)
			accountRepository.AssertNumberOfCalls(t, "GetById", tt.expectedAccountTypeRepositoryCall)
		})
	}
}

func TestTransactionService_CreateWithNegativeAmount(t *testing.T) {
	ctx := context.Background()
	t.Run("create transaction with negative amount when is a debit transaction", func(t *testing.T) {
		transactionRepository := new(TransactionMock)
		// validate here that the amount was transformed to negative
		transactionRepository.On("Create", ctx, model.Transaction{Amount: -1, Account: model.Account{Id: 1}, OperationTypes: model.OperationType{Id: 1}}).Return(model.Transaction{}, nil)

		operationTypeRepository := new(OperationTypeMock)
		operationTypeRepository.On("GetById", ctx, 1).Return(model.OperationType{Id: 1, Direction: model.DirectionDebit}, nil)
		//transaction with positive amount
		transaction := model.Transaction{Amount: 1, Account: model.Account{Id: 1}, OperationTypes: model.OperationType{Id: 1}}

		accountRepository := new(AccountRepositoryMock)
		accountRepository.On("GetById", ctx, 1).Return(model.Account{Id: 1}, nil)

		_, err := NewTransactionService(transactionRepository, operationTypeRepository, accountRepository).Create(ctx, transaction)
		assert.Nil(t, err)

		transactionRepository.AssertNumberOfCalls(t, "Create", 1)
		operationTypeRepository.AssertNumberOfCalls(t, "GetById", 1)
		accountRepository.AssertNumberOfCalls(t, "GetById", 1)
	})
}

type TransactionMock struct {
	mock.Mock
}

func (m *TransactionMock) Create(ctx context.Context, transaction model.Transaction) (model.Transaction, error) {
	args := m.Called(ctx, transaction)
	return args.Get(0).(model.Transaction), args.Error(1)
}

var _ repository.Transaction = (*TransactionMock)(nil)

type OperationTypeMock struct {
	mock.Mock
}

func (m *OperationTypeMock) GetById(ctx context.Context, id int) (model.OperationType, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.OperationType), args.Error(1)
}

var _ repository.OperationType = (*OperationTypeMock)(nil)
