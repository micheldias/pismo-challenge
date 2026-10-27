package service

import (
	"context"
	"pismo-challenge/internal/model"
	"pismo-challenge/internal/repository"
)

type Account interface {
	Create(context.Context, model.Account) (model.Account, error)
	GetById(context.Context, int) (model.Account, error)
}

type accountService struct {
	accountRepository repository.Account
}

func (a accountService) Create(ctx context.Context, account model.Account) (model.Account, error) {
	if err := account.Validate(); err != nil {
		return model.Account{}, err
	}
	return a.accountRepository.Create(ctx, account)
}

func (a accountService) GetById(ctx context.Context, id int) (model.Account, error) {
	return a.accountRepository.GetById(ctx, id)
}

// NewAccountService create an instance of account service
func NewAccountService(accountRepository repository.Account) Account {
	return &accountService{
		accountRepository: accountRepository,
	}
}

// --- INTERFACE GUARD ---
var _ Account = (*accountService)(nil)
