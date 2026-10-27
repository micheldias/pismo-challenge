package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	contextapi "pismo-challenge/internal/infra/context"
	"pismo-challenge/internal/infra/logging"
	servermiddleware "pismo-challenge/internal/infra/server/middleware"
	servermodel "pismo-challenge/internal/infra/server/model"
	"pismo-challenge/internal/infra/settings"
	"pismo-challenge/internal/model"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func CreateWebServer(routes []servermodel.Route, tasks []servermodel.ShutdownTask) WebServer {
	r := chi.NewRouter()
	r.Use(servermiddleware.InitRequest)
	r.Use(servermiddleware.Logger)
	r.Use(servermiddleware.Recover)
	r.Use(servermiddleware.ConfigureHeadersInResponse)

	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Get("/health", func(w http.ResponseWriter, request *http.Request) {
		_, _ = w.Write([]byte(`{"status": "UP"}`))
	})
	for _, route := range routes {
		r.Method(route.Method.String(), route.Path, CustomHandler(route.HandlerFunc))
	}
	return WebServer{
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%s", settings.Settings.Port),
			Handler: http.TimeoutHandler(r, settings.Settings.TimeoutEndpoints, ""),
		},
		shutdownTasks: tasks,
	}
}

type WebServer struct {
	httpServer    *http.Server
	shutdownTasks []servermodel.ShutdownTask
}

func (w WebServer) Start() {
	go func() {
		slog.Info(fmt.Sprintf("Starting server on : %s", settings.Settings.Port))

		if err := w.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Could not start server:", logging.ErrorMessageFieldLog, err)
			os.Exit(1)
		}
	}()
	w.gracefulShutdown()
}

func (w WebServer) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), settings.Settings.ServerShutdownTimeout)
	defer cancel()
	if err := w.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("error on shutdown http server: %v", err)
	}

	for _, task := range w.shutdownTasks {
		if task.Task == nil {
			continue
		}
		if err := task.Task(ctx); err != nil {
			return fmt.Errorf("error on shutdown %s: %v", task.Name, err)
		}
	}

	return nil
}

func (w WebServer) gracefulShutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-stop

	slog.Info("Shutting down server...")
	if err := w.Shutdown(); err != nil {
		slog.Error("Server forced to shutdown.", logging.ErrorMessageFieldLog, err)
		os.Exit(1)
	}

	slog.Info("Server exited gracefully")
}

func CustomHandler(handler func(*http.Request) (*servermodel.CustomResponse, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := contextapi.GetLogger(ctx)

		var requestBody string
		if reqDump, _ := httputil.DumpRequest(r, true); reqDump != nil {
			requestBody = string(reqDump)
		}
		logger.Info("Incoming HTTP request", logging.HttpRequestBodyContentFieldLog, requestBody)

		response, err := handler(r)
		if err != nil {
			response = errorHandler(err)
		}
		if statusText := http.StatusText(response.Status); statusText == "" {
			response.Status = http.StatusOK
		}
		w.WriteHeader(response.Status)

		if response.Body != nil {
			if encodingErr := json.NewEncoder(w).Encode(response.Body); encodingErr != nil {
				logger.Error("failed to encode response body", logging.ErrorMessageFieldLog, encodingErr)
				if err := json.NewEncoder(w).Encode(servermodel.DefaultInternalServerErrorBody()); err != nil {
					logger.Error("failed to encode response", logging.ErrorMessageFieldLog, encodingErr)
				}
			}
		}

		var responseBody string

		if jsonBytes, err := json.Marshal(response.Body); err != nil {
			logger.Error("failed to marshal response body", logging.ErrorMessageFieldLog, err)
		} else {
			responseBody = string(jsonBytes)
		}
		var requestDuration int64
		if startTime, ok := contextapi.GetStartTime(ctx); ok {
			requestDuration = time.Since(startTime).Milliseconds()
		}

		logger.Info("Finished handling request", logging.HttpResponseStatusCodeFieldLog, response.Status, logging.HttpResponseBodyContentFieldLog, responseBody,
			logging.ErrorMessageFieldLog, err, logging.EventDurationFieldLog, requestDuration)
	}
}

func errorHandler(err error) *servermodel.CustomResponse {
	switch e := err.(type) {
	case model.EntityNorFoundError:
		return &servermodel.CustomResponse{
			Status: http.StatusNetworkAuthenticationRequired,
			Body: []servermodel.ErrorDto{{
				Code:    model.EntityNorFoundErrorCode,
				Message: e.Reason},
			},
		}
	case model.MissingRequiredInformation:
		return &servermodel.CustomResponse{
			Status: http.StatusBadRequest,
			Body: func() []servermodel.ErrorDto {
				var missingRequiredErrors []servermodel.ErrorDto
				for _, missingRequiredError := range strings.Split(e.Error(), ";") {
					missingRequiredErrors = append(missingRequiredErrors, servermodel.ErrorDto{
						Code:    model.BadRequestErrorCode,
						Message: missingRequiredError,
					})
				}
				return missingRequiredErrors
			}(),
		}
	default:
		if errors.Is(err, context.DeadlineExceeded) {
			return &servermodel.CustomResponse{
				Status: http.StatusInternalServerError,
			}
		}
		return &servermodel.CustomResponse{
			Status: http.StatusInternalServerError,
			Body:   servermodel.DefaultInternalServerErrorBody(),
		}
	}
}
