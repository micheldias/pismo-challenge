package service

import (
	"context"
	"pismo-challenge/internal/model"
	"pismo-challenge/internal/repository"
)

type Transaction interface {
	Create(ctx context.Context, transaction model.Transaction) (model.Transaction, error)
}

type transactionService struct {
	transactionRepository   repository.Transaction
	operationTypeRepository repository.OperationType
	accountRepository       repository.Account
}

// Create a valid transaction in the database
func (t transactionService) Create(ctx context.Context, transaction model.Transaction) (model.Transaction, error) {
	operationType, err := t.operationTypeRepository.GetById(ctx, transaction.OperationTypes.Id)
	if err != nil {
		return model.Transaction{}, err
	}

	if !operationType.IsValidateOperationType() {
		return model.Transaction{}, model.MissingRequiredInformation{Reason: "invalid operation type"}
	}

	if _, err := t.accountRepository.GetById(ctx, transaction.Account.Id); err != nil {
		return model.Transaction{}, err
	}

	if err := transaction.Validate(); err != nil {
		return model.Transaction{}, err
	}

	transaction.Amount = operationType.CalculateAmountBasedOnDirection(transaction.Amount)

	return t.transactionRepository.Create(ctx, transaction)
}

// NewTransactionService create an instance of transaction service
func NewTransactionService(transactionRepository repository.Transaction, operationTypeRepository repository.OperationType, accountRepository repository.Account) Transaction {
	return &transactionService{
		transactionRepository:   transactionRepository,
		operationTypeRepository: operationTypeRepository,
		accountRepository:       accountRepository,
	}
}

// --- INTERFACE GUARD ---
var _ Transaction = (*transactionService)(nil)
