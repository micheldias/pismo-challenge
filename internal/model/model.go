package model

import (
	"fmt"
	"strings"
	"time"
)

const (
	emptyFieldMessage = "%s [%s] cannot be empty"
)

type Account struct {
	Id             int
	DocumentNumber string
}

func (a Account) Validate() error {
	if strings.TrimSpace(a.DocumentNumber) == "" {
		return MissingRequiredInformation{
			Reason: getAccountEmptyFieldErrorMessage("document_number"),
		}
	}
	return nil
}

type Transaction struct {
	Id             int
	Account        Account
	Amount         int
	OperationTypes OperationType
	EventDate      time.Time
}

func (t Transaction) Validate() error {
	var errors []error

	if t.Amount <= 0 {
		errors = append(errors, MissingRequiredInformation{
			Reason: "amount must be greater than zero",
		})
	}

	if t.Account.Id <= 0 {
		errors = append(errors, MissingRequiredInformation{
			Reason: "invalid account id",
		})
	}

	if t.OperationTypes.Id <= 0 {
		errors = append(errors, MissingRequiredInformation{
			Reason: "invalid operation type",
		})
	}

	if len(errors) > 0 {
		return MissingRequiredInformation{
			Errors: errors,
		}
	}
	return nil
}

const (
	DirectionCredit = "CREDIT"
	DirectionDebit  = "DEBIT"
)

type OperationType struct {
	Id          int
	Description string
	Direction   string
}

func (o OperationType) CalculateAmountBasedOnDirection(amount int) int {
	if o.Direction == DirectionCredit {
		return amount
	}
	return amount * -1
}

func (o OperationType) IsValidateOperationType() bool {
	if strings.EqualFold(o.Direction, DirectionCredit) || strings.EqualFold(o.Direction, DirectionDebit) {
		return true
	}
	return false
}

func getAccountEmptyFieldErrorMessage(name string) string {
	return getEmptyFieldErrorMessage("Account", name)
}

func getEmptyFieldErrorMessage(field, name string) string {
	return fmt.Sprintf(emptyFieldMessage, field, name)
}
