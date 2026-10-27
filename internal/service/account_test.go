package service

import (
	"context"
	"errors"
	"testing"

	"pismo-challenge/internal/model"
	"pismo-challenge/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAccountService_Create(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name                   string
		account                model.Account
		mockReturned           func() (model.Account, error)
		expectedResult         func() (model.Account, error)
		expectedRepositoryCall int
	}{
		{
			name:    "should forward the error received by the repository",
			account: model.Account{DocumentNumber: "123456789"},
			mockReturned: func() (model.Account, error) {
				return model.Account{}, errors.New("some kind of error")
			},
			expectedResult: func() (model.Account, error) {
				return model.Account{}, errors.New("some kind of error")
			},
			expectedRepositoryCall: 1,
		},
		{
			name:         "should return the error without call the repository if an invalid info was received",
			account:      model.Account{DocumentNumber: ""},
			mockReturned: func() (model.Account, error) { return model.Account{}, nil },
			expectedResult: func() (model.Account, error) {
				return model.Account{}, model.MissingRequiredInformation{Reason: "Account [document_number] cannot be empty"}
			},
			expectedRepositoryCall: 0,
		},
		{
			name:                   "should return the created account successfully",
			account:                model.Account{DocumentNumber: "1"},
			mockReturned:           func() (model.Account, error) { return model.Account{Id: 1, DocumentNumber: "1"}, nil },
			expectedResult:         func() (model.Account, error) { return model.Account{Id: 1, DocumentNumber: "1"}, nil },
			expectedRepositoryCall: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(AccountRepositoryMock)
			repo.On("Create", ctx, tt.account).Return(tt.mockReturned())
			account, err := NewAccountService(repo).Create(ctx, tt.account)
			expectedAccount, expectedError := tt.expectedResult()
			assert.Equal(t, expectedAccount, account)
			assert.Equal(t, expectedError, err)
			repo.AssertNumberOfCalls(t, "Create", tt.expectedRepositoryCall)
		})
	}
}

func TestAccountService_GetById(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name                   string
		id                     int
		mockReturned           func() (model.Account, error)
		expectedResult         func() (model.Account, error)
		expectedRepositoryCall int
	}{
		{
			name: "should forward the error received by the repository",
			id:   1,
			mockReturned: func() (model.Account, error) {
				return model.Account{}, errors.New("some kind of error")
			},
			expectedResult: func() (model.Account, error) {
				return model.Account{}, errors.New("some kind of error")
			},
			expectedRepositoryCall: 1,
		},
		{
			name:                   "should return the created account successfully",
			id:                     1,
			mockReturned:           func() (model.Account, error) { return model.Account{Id: 1, DocumentNumber: "1"}, nil },
			expectedResult:         func() (model.Account, error) { return model.Account{Id: 1, DocumentNumber: "1"}, nil },
			expectedRepositoryCall: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(AccountRepositoryMock)
			repo.On("GetById", ctx, tt.id).Return(tt.mockReturned())
			account, err := NewAccountService(repo).GetById(ctx, tt.id)
			expectedAccount, expectedError := tt.expectedResult()
			assert.Equal(t, expectedAccount, account)
			assert.Equal(t, expectedError, err)
			repo.AssertNumberOfCalls(t, "GetById", tt.expectedRepositoryCall)
		})
	}

}

type AccountRepositoryMock struct {
	mock.Mock
}

func (m *AccountRepositoryMock) Create(ctx context.Context, account model.Account) (model.Account, error) {
	args := m.Called(ctx, account)
	result, _ := args.Get(0).(model.Account)
	return result, args.Error(1)
}

func (m *AccountRepositoryMock) GetById(ctx context.Context, id int) (model.Account, error) {
	args := m.Called(ctx, id)
	result, _ := args.Get(0).(model.Account)
	return result, args.Error(1)
}

var _ repository.Account = (*AccountRepositoryMock)(nil)
