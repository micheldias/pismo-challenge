package servermiddleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
	contextapi "pismo-challenge/internal/infra/context"
	"pismo-challenge/internal/infra/logging"
	servermodel "pismo-challenge/internal/infra/server/model"
	"pismo-challenge/internal/infra/settings"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
)

const (
	requestIdHeader            = "X-Request-Id"
	contentTypeHeader          = "Content-Type"
	applicationJsonContentType = "application/json"
)

func InitRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		ctx := r.Context()
		ctx = contextapi.SetRequestId(ctx, uuid.New().String())
		ctx = contextapi.SetStartTime(ctx, startTime)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ConfigureHeadersInResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(requestIdHeader, contextapi.GetRequestId(r.Context()))
		w.Header().Add(contentTypeHeader, applicationJsonContentType)
		next.ServeHTTP(w, r)
	})
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		reqId := contextapi.GetRequestId(ctx)
		logger := slog.With(logging.TraceIdFieldLog, reqId, logging.HttpRequestMethodFieldLog, r.Method, logging.URLFieldLog, r.URL.String(), logging.AppNameLog, settings.Settings.AppName)
		r = r.WithContext(contextapi.SetLogger(ctx, logger))

		next.ServeHTTP(w, r)
	})
}

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)

				logger := contextapi.GetLogger(r.Context())
				logger.Error("a panic happened", logging.ErrorMessageFieldLog, err, logging.ErrorStackTraceFieldLog, string(debug.Stack()))

				if encodingErr := json.NewEncoder(w).Encode(servermodel.DefaultInternalServerErrorBody()); encodingErr != nil {
					logger.Error("failed to encode response", logging.ErrorMessageFieldLog, encodingErr)
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}
