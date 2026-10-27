package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	servermodel "pismo-challenge/internal/infra/server/model"
	"pismo-challenge/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomHandler(t *testing.T) {
	tests := []struct {
		name                 string
		handler              func(*http.Request) (*servermodel.CustomResponse, error)
		expectedResponseBody string
		expectedStatus       int
	}{
		{
			name: "should return the body and status get by custom response",
			handler: func(r *http.Request) (*servermodel.CustomResponse, error) {
				return &servermodel.CustomResponse{
					Status: http.StatusCreated,
					Body:   map[string]string{"message": "success"},
				}, nil
			},
			expectedStatus:       http.StatusCreated,
			expectedResponseBody: `{"message": "success"}`,
		},
		{
			name: "should return 200 status code if status custom response is 0",
			handler: func(r *http.Request) (*servermodel.CustomResponse, error) {
				return &servermodel.CustomResponse{
					Body: map[string]string{"message": "success"},
				}, nil
			},
			expectedStatus:       http.StatusOK,
			expectedResponseBody: `{"message": "success"}`,
		},
		{
			name: "should return internal server error status code and the standard error response payload",
			handler: func(r *http.Request) (*servermodel.CustomResponse, error) {
				return nil, errors.New("some error")
			},
			expectedStatus:       http.StatusInternalServerError,
			expectedResponseBody: `[{"code":"CODE-500", "message":"internal server error"}]`,
		},
		{
			name: "should return bad request error status code and the standard error response payload with all errors",
			handler: func(r *http.Request) (*servermodel.CustomResponse, error) {
				return nil, model.MissingRequiredInformation{
					Errors: []error{
						model.MissingRequiredInformation{Reason: "error_msg_1"},
						model.MissingRequiredInformation{Reason: "error_msg_2"},
					}}
			},
			expectedStatus:       http.StatusBadRequest,
			expectedResponseBody: `[{"code":"CODE-400", "message":"error_msg_1"},{"code":"CODE-400", "message":"error_msg_2"}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			handler := CustomHandler(tt.handler)
			handler.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/", http.NoBody))

			assert.Equal(t, tt.expectedStatus, response.Code)
			assert.JSONEq(t, tt.expectedResponseBody, response.Body.String())

		})
	}
}
