package controller

import (
	"encoding/json"
	"net/http"
	servermodel "pismo-challenge/internal/infra/server/model"
	"pismo-challenge/internal/model"
	"pismo-challenge/internal/service"
)

type Transaction struct {
	transactionService service.Transaction
}

func NewTransaction(transactionService service.Transaction) Transaction {
	return Transaction{transactionService: transactionService}
}

// Create handler
//
//	@Summary		transaction endpoint
//	@Description	create a transaction
//	@Tags			Transaction
//	@Accept			json
//	@Produce		json
//	@Param			request	body		TransactionRequestDto	true	"transaction"
//	@Success		201		{object}	TransactionResponseDto
//	@Failure		400		{array}		servermodel.ErrorDto
//	@Failure		500		{array}		servermodel.ErrorDto
//	@Failure		503
//
//	@Header			200,400,500	{string}	x-request-id	"request id"
//
//	@Router			/transactions [post]
func (t Transaction) Create(r *http.Request) (*servermodel.CustomResponse, error) {
	var transaction TransactionRequestDto
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		return nil, model.MissingRequiredInformation{Reason: "Invalid JSON payload"}
	}

	if err := transaction.Validate(); err != nil {
		return nil, err
	}

	transactionCreated, err := t.transactionService.Create(ctx, transaction.ToDomain())
	if err != nil {
		return nil, err
	}

	return &servermodel.CustomResponse{
		Status: http.StatusCreated,
		Body:   fromTransactionDomainToDto(transactionCreated),
	}, nil
}
