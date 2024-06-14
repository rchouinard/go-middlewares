package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rchouinard/middlewares"
	"github.com/stretchr/testify/assert"
)

func TestRequestID(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, ok := middlewares.GetRequestIDFromContext(ctx)

		assert.True(t, ok)
		assert.Len(t, id, 26)
	})

	handler := middlewares.RequestID(nextHandler)
	handler.ServeHTTP(rec, req)
}
