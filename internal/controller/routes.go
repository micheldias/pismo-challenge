package controller

import (
	"net/http"
	servermodel "pismo-challenge/internal/infra/server/model"
)

// GetRoutes returns all controllers endpoints and their associated methods and handler functions.
func GetRoutes(accountController Account, transactionController Transaction) []servermodel.Route {
	return []servermodel.Route{
		{
			Method:      http.MethodGet,
			Path:        "/accounts/{accountId}",
			HandlerFunc: accountController.GetById,
		},
		{
			Method:      http.MethodPost,
			Path:        "/accounts",
			HandlerFunc: accountController.Create,
		},
		{
			Method:      http.MethodPost,
			Path:        "/transactions",
			HandlerFunc: transactionController.Create,
		},
	}
}
