package servermodel

import (
	"context"
	"net/http"
	"pismo-challenge/internal/model"
)

type HttpMethod string

const (
	GET  HttpMethod = http.MethodGet
	POST HttpMethod = http.MethodPost
)

func (m HttpMethod) String() string {
	return string(m)
}

type Route struct {
	Method      HttpMethod
	Path        string
	HandlerFunc func(r *http.Request) (*CustomResponse, error)
}

type CustomResponse struct {
	Status int
	Body   any
}

type ErrorDto struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func DefaultInternalServerErrorBody() []ErrorDto {
	return []ErrorDto{
		{
			Code:    model.InternalServerErrorCode,
			Message: "internal server error",
		},
	}
}

type ShutdownTask struct {
	Name string
	Task ShutdownTaskFunc
}

type ShutdownTaskFunc func(ctx context.Context) error
