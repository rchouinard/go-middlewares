package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rchouinard/middlewares"
	"github.com/stretchr/testify/assert"
)

func TestResponseWriterBeforeWrite(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := middlewares.NewResponseWriter(rec)

	assert.Equal(t, 0, rw.Status())
	assert.Equal(t, "", rec.Body.String())
	assert.Equal(t, 0, rw.Size())
	assert.Equal(t, false, rw.Written())
}

func TestResponseWriterWriteString(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := middlewares.NewResponseWriter(rec)

	content := "Hello, World!"
	rw.Write([]byte(content))

	assert.Equal(t, http.StatusOK, rw.Status())
	assert.Equal(t, content, rec.Body.String())
	assert.Equal(t, len(content), rw.Size())
	assert.True(t, rw.Written())
}

func TestResponseWriterWriteHeader(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := middlewares.NewResponseWriter(rec)

	rw.WriteHeader(http.StatusNotFound)

	assert.Equal(t, http.StatusNotFound, rw.Status())
	assert.Equal(t, "", rec.Body.String())
	assert.Equal(t, 0, rw.Size())
	assert.True(t, rw.Written())
}
