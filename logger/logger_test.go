package logger_test

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rchouinard/go-middlewares/logger"
	"github.com/stretchr/testify/assert"
)

type loggerWriter struct {
	Content []byte
}

func (lw *loggerWriter) Write(in []byte) (int, error) {
	b := len(in)

	lw.Content = append(lw.Content, in...)

	return b, nil
}

func TestLogger(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "100.100.100.100:12345"
	req.Header.Set("Referer", "http://127.0.0.1:8000/")
	req.Header.Set("User-Agent", "golang/test")

	content := "Hello, World!"

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if l, ok := logger.GetLoggerFromRequest(r); ok {
			l.Info(content)
			w.WriteHeader(200)
			fmt.Fprintf(w, content)
		}

		return
	})

	accessWriter := &loggerWriter{}
	errorWriter := &loggerWriter{}

	accessLogger := slog.New(slog.NewJSONHandler(accessWriter, nil))
	errorLogger := slog.New(slog.NewJSONHandler(errorWriter, nil))

	mw := logger.New(errorLogger, accessLogger)(nextHandler)
	mw.ServeHTTP(rec, req)

	var accessJSON map[string]interface{}
	json.Unmarshal(accessWriter.Content, &accessJSON)
	assert.NotEmpty(t, accessJSON["time"].(string))
	assert.Equal(t, "100.100.100.100", accessJSON["remote_ip"].(string))
	assert.Equal(t, "GET / HTTP/1.1", accessJSON["request"].(string))
	assert.Equal(t, 200, int(accessJSON["response"].(float64)))
	assert.Equal(t, len(content), int(accessJSON["bytes"].(float64)))
	assert.Equal(t, "http://127.0.0.1:8000/", accessJSON["referer"].(string))
	assert.Equal(t, "golang/test", accessJSON["agent"].(string))
	assert.Greater(t, int(accessJSON["duration"].(float64)), 1)

	var errorJSON map[string]interface{}
	json.Unmarshal(errorWriter.Content, &errorJSON)
	assert.NotEmpty(t, errorJSON["time"].(string))
	assert.Equal(t, "INFO", errorJSON["level"].(string))
	assert.Equal(t, content, errorJSON["msg"].(string))
}
