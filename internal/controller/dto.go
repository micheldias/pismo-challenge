package controller

import (
	"pismo-challenge/internal/model"
)

const DateTimeFormat = "2006-01-02 15:04:05Z"

// TransactionRequestDto struct represent the transaction requests
type TransactionRequestDto struct {
	//account identifier
	AccountId int `json:"account_id" validate:"required,gt=0" example:"1"`
	//operation type id
	OperationTypeId int `json:"operation_type_id" validate:"required,gt=0" example:"1"`
	//transaction amount in cents
	Amount int `json:"amount" validate:"required,gt=0" example:"1"`
}

func (t TransactionRequestDto) ToDomain() model.Transaction {
	return model.Transaction{
		Account:        model.Account{Id: t.AccountId},
		Amount:         t.Amount,
		OperationTypes: model.OperationType{Id: t.OperationTypeId},
	}
}

func (t TransactionRequestDto) Validate() error {
	var errors []error

	if t.Amount <= 0 {
		errors = append(errors, model.MissingRequiredInformation{
			Reason: "amount must be greater than zero",
		})
	}

	if t.AccountId <= 0 {
		errors = append(errors, model.MissingRequiredInformation{
			Reason: "invalid account id",
		})
	}

	if t.OperationTypeId <= 0 {
		errors = append(errors, model.MissingRequiredInformation{
			Reason: "invalid operation type",
		})
	}

	if len(errors) > 0 {
		return model.MissingRequiredInformation{
			Errors: errors,
		}
	}
	return nil
}

// TransactionResponseDto struct represent the transaction response
type TransactionResponseDto struct {
	//transaction identifier
	TransactionID int `json:"transaction_id"`
	//account identifier
	AccountId int `json:"account_id"`
	//operation type id
	OperationTypeId int `json:"operation_type_id"`
	//transaction amount in cents
	Amount int `json:"amount"`
	//transaction created date
	EventDate string `json:"event_date"`
}

func fromTransactionDomainToDto(transaction model.Transaction) TransactionResponseDto {
	return TransactionResponseDto{
		TransactionID:   transaction.Id,
		AccountId:       transaction.Account.Id,
		OperationTypeId: transaction.OperationTypes.Id,
		Amount:          transaction.Amount,
		EventDate:       transaction.EventDate.Format(DateTimeFormat),
	}
}

// AccountResponseDto struct represent the account response
type AccountResponseDto struct {
	//account identifier
	Id int `json:"account_id"`
	//customer document number
	DocumentNumber string `json:"document_number"`
}

func fromAccountDomainToDto(account model.Account) AccountResponseDto {
	return AccountResponseDto{
		Id:             account.Id,
		DocumentNumber: account.DocumentNumber,
	}
}

// AccountRequestDto struct represent the account request
type AccountRequestDto struct {
	//customer document number
	DocumentNumber string `json:"document_number" validate:"required,alphanum,min=4,max=100" example:"ABC1234"`
}

func (a AccountRequestDto) Validate() error {
	if !isAlphanumeric(a.DocumentNumber) {
		return model.MissingRequiredInformation{
			Reason: "document number must be alphanumeric",
		}
	}

	return nil
}

func (a AccountRequestDto) ToDomain() model.Account {
	return model.Account{
		DocumentNumber: a.DocumentNumber,
	}
}

func isAlphanumeric(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')) {
			return false
		}
	}
	return true
}
