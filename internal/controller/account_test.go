package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	servermodel "pismo-challenge/internal/infra/server/model"
	"pismo-challenge/internal/model"
	"pismo-challenge/internal/service"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAccountHandler_GetById(t *testing.T) {
	tests := []struct {
		name           string
		mockReturned   func() (model.Account, error)
		id             string
		expectedAssert func(*testing.T, *servermodel.CustomResponse, error)
	}{
		{
			name: "validate success response",
			id:   "123",
			mockReturned: func() (model.Account, error) {
				return model.Account{Id: 1, DocumentNumber: "123456"}, nil
			},
			expectedAssert: func(t *testing.T, response *servermodel.CustomResponse, err error) {
				expectedResponseBody := `{"account_id":1,"document_number":"123456"}`
				responseBody, errForMarshal := json.Marshal(response.Body)
				assert.Nil(t, errForMarshal)
				assert.Nil(t, err)
				assert.JSONEq(t, expectedResponseBody, string(responseBody))
				assert.Equal(t, http.StatusOK, response.Status)
			},
		},
		{
			name: "should forward the the error received by the service",
			id:   "123",
			mockReturned: func() (model.Account, error) {
				return model.Account{}, errors.New("some error")
			},
			expectedAssert: func(t *testing.T, response *servermodel.CustomResponse, err error) {
				assert.Nil(t, response)
				assert.EqualError(t, err, "some error")
			},
		},
		{
			name:         "should return error if account id is not a number",
			id:           "123a",
			mockReturned: func() (model.Account, error) { return model.Account{}, nil },
			expectedAssert: func(t *testing.T, response *servermodel.CustomResponse, err error) {
				assert.Nil(t, response)
				assert.Equal(t, model.MissingRequiredInformation{Reason: "account id is not a number greater than zero"}, err)
			},
		},
		{
			name:         "should return error if account id is zero",
			id:           "0",
			mockReturned: func() (model.Account, error) { return model.Account{}, nil },
			expectedAssert: func(t *testing.T, response *servermodel.CustomResponse, err error) {
				assert.Nil(t, response)
				assert.Equal(t, model.MissingRequiredInformation{Reason: "account id is not a number greater than zero"}, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountService := new(AccountServiceMock)

			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("accountId", tt.id)

			req := httptest.NewRequest(http.MethodGet, "/accounts/{accountId}", http.NoBody)
			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx)
			req = req.WithContext(ctx)

			accountService.On("GetById", ctx, 123).Return(tt.mockReturned())
			handler := NewAccount(accountService)

			response, err := handler.GetById(req)
			tt.expectedAssert(t, response, err)

		})
	}

}

func TestAccountHandler_Create(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		mockReturned   func() (model.Account, error)
		id             string
		expectedAssert func(*testing.T, *servermodel.CustomResponse, error)
	}{
		{
			name:        "validate success response",
			requestBody: `{"document_number":"123456"}`,
			mockReturned: func() (model.Account, error) {
				return model.Account{Id: 1, DocumentNumber: "123456"}, nil
			},
			expectedAssert: func(t *testing.T, response *servermodel.CustomResponse, err error) {
				expectedResponseBody := `{"account_id":1,"document_number":"123456"}`
				responseBody, errForMarshal := json.Marshal(response.Body)
				assert.Nil(t, errForMarshal)
				assert.Nil(t, err)
				assert.JSONEq(t, expectedResponseBody, string(responseBody))
				assert.Equal(t, http.StatusCreated, response.Status)
			},
		},
		{
			name:        "return error if service return error",
			requestBody: `{"document_number":"123456"}`,
			mockReturned: func() (model.Account, error) {
				return model.Account{}, errors.New("some error")
			},
			expectedAssert: func(t *testing.T, response *servermodel.CustomResponse, err error) {
				assert.Nil(t, response)
				assert.EqualError(t, err, "some error")
			},
		},
		{
			name:         "return error if is an invalid body",
			requestBody:  `{"document_num`,
			mockReturned: func() (model.Account, error) { return model.Account{}, nil },
			expectedAssert: func(t *testing.T, response *servermodel.CustomResponse, err error) {
				assert.Nil(t, response)
				assert.Equal(t, model.MissingRequiredInformation{Reason: "Invalid JSON payload"}, err)
			},
		},
		{
			name:         "return error if the document is empty",
			requestBody:  `{"document_number":""}`,
			mockReturned: func() (model.Account, error) { return model.Account{}, nil },
			expectedAssert: func(t *testing.T, response *servermodel.CustomResponse, err error) {
				assert.Nil(t, response)
				assert.Equal(t, model.MissingRequiredInformation{Reason: "document number must be alphanumeric"}, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountService := new(AccountServiceMock)
			accountService.On("Create", context.Background(), model.Account{DocumentNumber: "123456"}).Return(tt.mockReturned())

			req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer([]byte(tt.requestBody)))
			response, err := NewAccount(accountService).Create(req)
			tt.expectedAssert(t, response, err)

		})
	}

}

type AccountServiceMock struct {
	mock.Mock
}

func (m *AccountServiceMock) Create(ctx context.Context, account model.Account) (model.Account, error) {
	args := m.Called(ctx, account)
	return args.Get(0).(model.Account), args.Error(1)
}

func (m *AccountServiceMock) GetById(ctx context.Context, id int) (model.Account, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Account), args.Error(1)
}

var _ service.Account = (*AccountServiceMock)(nil)
