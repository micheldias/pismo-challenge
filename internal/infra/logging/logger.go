package logging

import (
	"io"
	"log/slog"
)

func ConfigureStructuredLog(writer io.Writer) {
	logger := slog.New(slog.NewJSONHandler(writer, &slog.HandlerOptions{}))
	slog.SetDefault(logger)
}

const (
	AppNameLog                      = "app.name"
	EventDurationFieldLog           = "request.duration"
	ErrorMessageFieldLog            = "error.message"
	HttpRequestBodyContentFieldLog  = "request.body"
	HttpResponseStatusCodeFieldLog  = "response.status_code"
	HttpResponseBodyContentFieldLog = "response.body"
	ErrorStackTraceFieldLog         = "stack_trace"
	HttpRequestMethodFieldLog       = "request.method"
	URLFieldLog                     = "path"
	TraceIdFieldLog                 = "trace.id"
)
