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
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransactionHandler_Create(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		mockReturned   func() (model.Transaction, error)
		expectedAssert func(*testing.T, *servermodel.CustomResponse, error)
	}{
		{
			name:        "validate success response",
			requestBody: `{"account_id":1,"operation_type_id":1,"amount":1000}`,
			mockReturned: func() (model.Transaction, error) {
				return model.Transaction{
					Id:             1,
					Account:        model.Account{Id: 1, DocumentNumber: "123456"},
					Amount:         1000,
					OperationTypes: model.OperationType{Id: 1},
					EventDate:      time.Date(2009, time.November, 17, 20, 34, 58, 0, time.UTC),
				}, nil
			},
			expectedAssert: func(t *testing.T, response *servermodel.CustomResponse, err error) {
				expectedResponseBody := `{"transaction_id":1,"account_id":1,"operation_type_id":1,"amount":1000,"event_date":"2009-11-17 20:34:58Z"}`
				responseBody, errForMarshal := json.Marshal(response.Body)
				assert.Nil(t, errForMarshal)
				assert.Nil(t, err)
				assert.JSONEq(t, expectedResponseBody, string(responseBody))
				assert.Equal(t, http.StatusCreated, response.Status)
			},
		},
		{
			name:        "return error if service return error",
			requestBody: `{"account_id":1,"operation_type_id":1,"amount":1000}`,
			mockReturned: func() (model.Transaction, error) {
				return model.Transaction{}, errors.New("some error")
			},
			expectedAssert: func(t *testing.T, response *servermodel.CustomResponse, err error) {
				assert.Nil(t, response)
				assert.EqualError(t, err, "some error")
			},
		},
		{
			name:         "return error if is an invalid body",
			requestBody:  `{"document_num`,
			mockReturned: func() (model.Transaction, error) { return model.Transaction{}, nil },
			expectedAssert: func(t *testing.T, response *servermodel.CustomResponse, err error) {
				assert.Nil(t, response)
				assert.Equal(t, model.MissingRequiredInformation{Reason: "Invalid JSON payload"}, err)
			},
		},
		{
			name:        "return error if is body contains invalid information",
			requestBody: `{"account_id":0,"operation_type_id":0,"amount":0}`,
			mockReturned: func() (model.Transaction, error) {
				return model.Transaction{}, errors.New("some error")
			},
			expectedAssert: func(t *testing.T, response *servermodel.CustomResponse, err error) {
				assert.Nil(t, response)
				assert.EqualError(t, err, "amount must be greater than zero;invalid account id;invalid operation type")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transactionService := new(ServiceTransactionMock)
			transactionService.On("Create", context.Background(), model.Transaction{Account: model.Account{Id: 1}, Amount: 1000, OperationTypes: model.OperationType{Id: 1}}).Return(tt.mockReturned())

			req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer([]byte(tt.requestBody)))
			response, err := NewTransaction(transactionService).Create(req)
			tt.expectedAssert(t, response, err)
		})
	}

}

type ServiceTransactionMock struct {
	mock.Mock
}

func (m *ServiceTransactionMock) Create(ctx context.Context, transaction model.Transaction) (model.Transaction, error) {
	args := m.Called(ctx, transaction)
	return args.Get(0).(model.Transaction), args.Error(1)
}

var _ service.Transaction = (*ServiceTransactionMock)(nil)
