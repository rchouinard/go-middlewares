package middlewares

import (
	"context"
	"net/http"

	"github.com/oklog/ulid/v2"
)

const (
	requestIDKey ctxKey = "requestID"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := ulid.Make().String()
		ctx := context.WithValue(r.Context(), requestIDKey, id)
		w.Header().Add("X-Request-Id", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetRequestIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(requestIDKey).(string)
	return id, ok
}

func GetRequestIDFromRequest(r *http.Request) (string, bool) {
	ctx := r.Context()
	return GetRequestIDFromContext(ctx)
}
