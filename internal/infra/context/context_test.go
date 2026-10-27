package contextapi

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetLogger(t *testing.T) {
	t.Run("set logger with nil value, should add a default logger in context", func(t *testing.T) {
		ctx := SetLogger(context.Background(), nil)
		assert.NotNil(t, ctx.Value(loggerKey).(*slog.Logger))
	})
}

func TestGetLogger(t *testing.T) {
	t.Run("get logger with context without logger, should return a default logger", func(t *testing.T) {
		assert.NotNil(t, GetLogger(context.Background()))
	})
}

func TestSetStartTime(t *testing.T) {
	t.Run("should set value in context successfully", func(t *testing.T) {
		startTime := time.Now()
		ctxWithStartTime := SetStartTime(context.Background(), startTime)

		gotTime := ctxWithStartTime.Value(startTimeKey).(time.Time)

		assert.Equal(t, startTime, gotTime)
	})

	t.Run("dont set start time if context is nil", func(t *testing.T) {
		assert.Nil(t, SetStartTime(nil, time.Now()))

	})
}

func TestGetStartTime(t *testing.T) {
	t.Run("get default if start time  doesn't exist in context", func(t *testing.T) {
		gotTime, ok := GetStartTime(context.Background())
		assert.False(t, ok)
		assert.Equal(t, time.Time{}, gotTime)
	})
	t.Run("get start time successfully if exists context", func(t *testing.T) {
		startTime := time.Now()
		ctx := context.WithValue(context.Background(), startTimeKey, startTime)

		gotTime, ok := GetStartTime(ctx)
		assert.True(t, ok)
		assert.Equal(t, startTime, gotTime)
	})
}

func TestGetRequestId(t *testing.T) {
	t.Run("get empty string value doesn't exist in context", func(t *testing.T) {
		assert.Empty(t, GetRequestId(context.Background()))
	})
	t.Run("get value if exists in context", func(t *testing.T) {
		requestId := "request-id"
		ctx := context.WithValue(context.Background(), requestIdKey, requestId)
		assert.Equal(t, "request-id", GetRequestId(ctx))
	})
}

func TestSetRequestId(t *testing.T) {
	t.Run("should set value in context successfully", func(t *testing.T) {
		requestId := "request-id"
		ctx := SetRequestId(context.Background(), requestId)
		assert.Equal(t, requestId, ctx.Value(requestIdKey).(string))
	})

	t.Run("dont set request id if context is nil", func(t *testing.T) {
		assert.Nil(t, SetRequestId(nil, "request-id"))

	})
}
