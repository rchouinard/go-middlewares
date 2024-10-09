package requestid

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
)

const (
	requestIDKey string = "requestID"
)

// IDGenerator is a type definition for a function that generates a unique ID as a string.
type IDGenerator func() string

// New is a higher-order function that takes an IDGenerator function as an argument
// and returns a function that takes an http.Handler as an argument and returns an http.Handler.
//
// Parameters:
// - genFn: A function that generates a unique ID.
//
// Returns:
// - func(http.Handler) http.Handler: A function that wraps the provided http.Handler with request ID generation and setting.
func New(genFn IDGenerator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := genFn()
			ctx := context.WithValue(r.Context(), requestIDKey, id)
			w.Header().Add("X-Request-Id", id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequestID is a middleware function that wraps the provided http.Handler with request ID generation and setting.
// It also adds a "X-Request-Id" header to the response.
//
// Parameters:
// - next: The http.Handler to be wrapped with request ID generation and setting.
//
// Returns:
// - http.Handler: The wrapped http.Handler.
func RequestID(next http.Handler) http.Handler {
	return New(ULIDGenerator)(next)
}

// GetRequestIDFromContext extracts the request ID from the provided context.
//
// Parameters:
// - ctx: The context from which to extract the request ID.
//
// Returns:
// - string: The extracted request ID.
// - bool: A boolean value indicating whether the request ID was successfully extracted.
func GetRequestIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(requestIDKey).(string)
	return id, ok
}

// GetRequestIDFromRequest extracts the request ID from the provided HTTP request's context.
//
// Parameters:
// - r: A pointer to the HTTP request from which to extract the request ID.
//
// Returns:
// - string: The extracted request ID.
// - bool: A boolean value indicating whether the request ID was successfully extracted.
func GetRequestIDFromRequest(r *http.Request) (string, bool) {
	ctx := r.Context()
	return GetRequestIDFromContext(ctx)
}

// ULIDGenerator is a function that generates a new ULID as a string.
func ULIDGenerator() string {
	return ulid.Make().String()
}

// UUIDGenerator is a function that generates a new UUID as a string.
func UUIDGenerator() string {
	return uuid.NewString()
}
