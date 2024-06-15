package middlewares

import (
	"context"
	"net/http"

	"github.com/oklog/ulid/v2"
)

const (
	requestIDKey ctxKey = "requestID"
)

// RequestID implements a middleware which adds a unique ID to each request.
//
// The request ID can be retrieved from the X-Request-Id response header or from the request context with the
// [GetRequestIDFromContext] or [GetRequestIDFromRequest] functions.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := ulid.Make().String()
		ctx := context.WithValue(r.Context(), requestIDKey, id)
		w.Header().Add("X-Request-Id", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRequestIDFromContext returns the unique ID from a given request context.
func GetRequestIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(requestIDKey).(string)
	return id, ok
}

// GetRequestIDFromContext returns the unique ID from a given request.
func GetRequestIDFromRequest(r *http.Request) (string, bool) {
	ctx := r.Context()
	return GetRequestIDFromContext(ctx)
}
