package servermiddleware

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	contextapi "pismo-challenge/internal/infra/context"
	"pismo-challenge/internal/infra/logging"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitRequest(t *testing.T) {
	t.Run("should generate requestId and startTime in context", func(t *testing.T) {
		InitRequest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := contextapi.GetRequestId(r.Context())
			assert.NotEmpty(t, reqID)
			_, startTimeGenerated := contextapi.GetStartTime(r.Context())
			assert.True(t, startTimeGenerated)

		})).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/test", http.NoBody))
	})

}

func TestLogger(t *testing.T) {
	t.Run("should contain the requestId in logs", func(t *testing.T) {
		sink := &bytes.Buffer{}
		logging.ConfigureStructuredLog(sink)

		ctx := contextapi.SetRequestId(context.Background(), "request-id-context")
		Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := contextapi.GetLogger(r.Context())
			logger.Info("test")
		})).ServeHTTP(httptest.NewRecorder(), httptest.NewRequestWithContext(ctx, http.MethodGet, "/hello", http.NoBody))
		logFields := readLogFields(t, sink)
		assert.Equal(t, "request-id-context", logFields[logging.TraceIdFieldLog].(string))
	})

}

func TestConfigureHeadersInResponse(t *testing.T) {
	t.Run("should set headers in response", func(t *testing.T) {
		response := httptest.NewRecorder()
		ctx := contextapi.SetRequestId(context.Background(), "request-id-context")
		ConfigureHeadersInResponse(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		})).ServeHTTP(response, httptest.NewRequestWithContext(ctx, http.MethodGet, "/hello", http.NoBody))
		responseHeaders := response.Header()
		assert.Equal(t, "request-id-context", responseHeaders[requestIdHeader][0])
		assert.Equal(t, applicationJsonContentType, responseHeaders[contentTypeHeader][0])
	})

}

func TestRecover(t *testing.T) {
	t.Run("should log panic error, stack trace and return internal server error in response", func(t *testing.T) {
		sink := &bytes.Buffer{}
		logging.ConfigureStructuredLog(sink)
		response := httptest.NewRecorder()
		Recover(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("panic error")
		})).ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/hello", http.NoBody))
		logFields := readLogFields(t, sink)
		assert.NotEmpty(t, logFields[logging.ErrorStackTraceFieldLog].(string))
		assert.Equal(t, "a panic happened", logFields["msg"].(string))
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.JSONEq(t, `[{"code":"CODE-500", "message":"internal server error"}]`, response.Body.String())

	})

}

func readLogFields(t *testing.T, sink *bytes.Buffer) map[string]any {
	scanner := bufio.NewScanner(sink)
	scanner.Scan()

	var logFields map[string]any
	if err := json.Unmarshal(scanner.Bytes(), &logFields); err != nil {
		t.Fatalf("invalid json format: %v", err)
	}
	return logFields
}
