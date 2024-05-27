package middlewares_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rchouinard/middlewares"
	"github.com/stretchr/testify/assert"
)

type loggerReaderWriter struct {
	Bytes int
}

func (lrw *loggerReaderWriter) Write(in []byte) (int, error) {
	b := len(in)

	lrw.Bytes += b

	return b, nil
}

func TestLogger(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if el, ok := middlewares.LoggerFromRequest(r); ok {
			el.Info("This is a test message")
		}

		return
	})

	accessWriter := &loggerReaderWriter{}
	errorWriter := &loggerReaderWriter{}

	assert.Equal(t, 0, accessWriter.Bytes)
	assert.Equal(t, 0, errorWriter.Bytes)

	handler := middlewares.LoggerWithConfig(middlewares.LoggerConfig{
		AccessHandler: slog.NewJSONHandler(accessWriter, nil),
		ErrorHandler:  slog.NewJSONHandler(errorWriter, nil),
	})(nextHandler)
	handler.ServeHTTP(rec, req)

	assert.Greater(t, accessWriter.Bytes, 0)
	assert.Greater(t, errorWriter.Bytes, 0)
}
