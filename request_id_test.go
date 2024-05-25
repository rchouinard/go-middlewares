package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rchouinard/middlewares"
	"github.com/stretchr/testify/assert"
)

func TestRequestIDMiddleware(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, ok := middlewares.RequestIDFromContext(ctx)

		assert.True(t, ok)
		assert.Len(t, id, 26)
	})

	handler := middlewares.RequestIDMiddleware(nextHandler)
	handler.ServeHTTP(rec, req)
}
