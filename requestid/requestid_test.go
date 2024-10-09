package requestid_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rchouinard/go-middlewares/requestid"
	"github.com/stretchr/testify/assert"
)

func TestNewRequestID(t *testing.T) {
	new1 := requestid.New(DummyGenerator)
	assert.IsType(t, requestid.RequestID, new1)
}

func TestULIDRequestID(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, ok := requestid.GetRequestIDFromContext(ctx)

		assert.True(t, ok)
		assert.Len(t, id, 26)
	})

	handler := requestid.New(requestid.ULIDGenerator)(nextHandler)
	handler.ServeHTTP(rec, req)
}

func TestUUIDRequestID(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, ok := requestid.GetRequestIDFromContext(ctx)

		assert.True(t, ok)
		assert.Len(t, id, 36)
	})

	handler := requestid.New(requestid.UUIDGenerator)(nextHandler)
	handler.ServeHTTP(rec, req)
}

func DummyGenerator() string {
	return "dummy"
}
