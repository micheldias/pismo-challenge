package contextapi

import (
	"context"
	"log/slog"
	"time"
)

type contextKey int

const (
	startTimeKey contextKey = iota + 1
	loggerKey
	requestIdKey
)

func SetLogger(ctx context.Context, logger *slog.Logger) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if logger == nil {
		logger = slog.Default()
	}

	return context.WithValue(ctx, loggerKey, logger)
}

func GetLogger(ctx context.Context) *slog.Logger {
	if ctx != nil {
		if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
			return logger
		}
	}
	return slog.Default()

}

func GetStartTime(ctx context.Context) (time.Time, bool) {
	if ctx != nil {
		startTime, ok := ctx.Value(startTimeKey).(time.Time)
		return startTime, ok
	}
	return time.Time{}, false
}

func SetStartTime(ctx context.Context, startTime time.Time) context.Context {
	if ctx != nil {
		return context.WithValue(ctx, startTimeKey, startTime)
	}
	return ctx
}

func GetRequestId(ctx context.Context) string {
	var requestId string
	if ctx != nil {
		requestId, _ = ctx.Value(requestIdKey).(string)
	}
	return requestId
}

func SetRequestId(ctx context.Context, requestId string) context.Context {
	if ctx != nil {
		return context.WithValue(ctx, requestIdKey, requestId)
	}
	return ctx
}
