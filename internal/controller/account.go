package controller

import (
	"encoding/json"
	"net/http"
	servermodel "pismo-challenge/internal/infra/server/model"
	"pismo-challenge/internal/model"
	"pismo-challenge/internal/service"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Account struct {
	accountService service.Account
}

func NewAccount(accountService service.Account) Account {
	return Account{accountService: accountService}
}

// Create handler
//
//	@Summary		account endpoint
//	@Description	create an account
//	@Tags			Account
//	@Accept			json
//	@Produce		json
//	@Param			request	body		AccountRequestDto	true	"account"
//	@Success		201		{object}	AccountResponseDto
//	@Failure		400		{array}		servermodel.ErrorDto
//	@Failure		500		{array}		servermodel.ErrorDto
//	@Failure		503
//
//	@Header			200,400,500	{string}	x-request-id	"request id"
//
//	@Router			/accounts [post]
func (a Account) Create(r *http.Request) (*servermodel.CustomResponse, error) {
	var accountDto AccountRequestDto
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&accountDto); err != nil {
		return nil, model.MissingRequiredInformation{Reason: "Invalid JSON payload"}
	}

	if err := accountDto.Validate(); err != nil {
		return nil, err
	}

	accountCreated, err := a.accountService.Create(ctx, accountDto.ToDomain())
	if err != nil {
		return nil, err
	}

	return &servermodel.CustomResponse{
		Status: http.StatusCreated,
		Body:   fromAccountDomainToDto(accountCreated),
	}, nil
}

// GetById handler
//
//	@Summary		get account by id
//	@Description	Get an account by account id
//	@Tags			Account
//	@Accept			json
//	@Produce		json
//	@Param			accountId	path		int	true	"Account ID"
//	@Success		201			{object}	AccountResponseDto
//	@Failure		400			{array}		servermodel.ErrorDto
//	@Failure		500			{array}		servermodel.ErrorDto
//	@Failure		503
//
//	@Header			200,400,500	{string}	x-request-id	"request id"
//
//	@Router			/accounts/{accountId} [get]
func (a Account) GetById(r *http.Request) (*servermodel.CustomResponse, error) {
	ctx := r.Context()
	accountIdParam := chi.URLParam(r, "accountId")
	accountId, err := strconv.Atoi(accountIdParam)

	if err != nil || accountId < 1 {
		return nil, model.MissingRequiredInformation{Reason: "account id is not a number greater than zero"}
	}

	account, err := a.accountService.GetById(ctx, accountId)
	if err != nil {
		return nil, err
	}

	return &servermodel.CustomResponse{
		Status: http.StatusOK,
		Body:   fromAccountDomainToDto(account),
	}, nil
}
