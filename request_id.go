package middlewares

import (
	"context"
	"net/http"

	"github.com/oklog/ulid/v2"
)

const (
	requestIDKey string = "requestID"
)

func RequestIDMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := ulid.Make().String()
		ctx := context.WithValue(r.Context(), requestIDKey, id)
		w.Header().Add("X-Request-Id", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequestIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(requestIDKey).(string)
	return id, ok
}

func RequestIDFromRequest(r *http.Request) (string, bool) {
	ctx := r.Context()
	return RequestIDFromContext(ctx)
}
