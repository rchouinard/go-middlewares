package logger

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/rchouinard/go-middlewares/common"
	"github.com/rchouinard/go-middlewares/requestid"
)

const (
	loggerKey string = "requestID"
)

var (
	logHandlerOpts = &slog.HandlerOptions{
		Level: slog.LevelError,
	}

	accLogHandlerOpts = &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.MessageKey || a.Key == slog.LevelKey {
				return slog.Attr{} // drop attribute
			}

			return a
		},
	}
)

// New returns a middleware that logs requests with the provided logger and access logger.
// It attaches the logger to the request context and logs request details using the access logger.
//
// Parameters:
// - logger: The logger used for error logging.
// - accLogger: The logger used for access logging.
//
// Returns:
// - func(http.Handler) http.Handler: A function that wraps the provided http.Handler with error and request logging.
func New(logger, accLogger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// create new request context with logger attached
			// this will be passed down to the remaining handlers
			ctx := context.WithValue(r.Context(), loggerKey, logger)

			// run the remaining handlers and gather stats
			lrw := common.NewResponseWriter(w)
			start := time.Now()
			next.ServeHTTP(lrw, r.WithContext(ctx))
			duration := time.Since(start)

			// drop the port from the remote address
			addr, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				addr = r.RemoteAddr
			}

			// log the request id, if available
			var attrs []slog.Attr
			if id, ok := requestid.GetRequestIDFromRequest(r); ok {
				attrs = append(attrs, slog.String("request_id", id))

			}

			f := append(attrs,
				slog.String("remote_ip", addr),
				slog.String("request", fmt.Sprintf("%s %s %s", r.Method, r.URL.String(), r.Proto)),
				slog.Int("response", lrw.Status()),
				slog.Int("bytes", lrw.Size()),
				slog.String("referer", r.Referer()),
				slog.String("agent", r.UserAgent()),
				slog.Duration("duration", duration),
			)

			accLogger.LogAttrs(context.Background(), slog.LevelInfo, "access_log", f...)
		})
	}
}

// Logger is a middleware that wraps the provided http.Handler with text error and request logging.
//
// Parameters:
// - next: The http.Handler to be wrapped with logging.
//
// Returns:
// - http.Handler: The wrapped http.Handler.
func Logger(next http.Handler) http.Handler {
	logger := slog.New(slog.NewTextHandler(os.Stderr, logHandlerOpts))
	accLogger := slog.New(slog.NewTextHandler(os.Stdout, accLogHandlerOpts))

	return New(logger, accLogger)(next)
}

// Logger is a middleware that wraps the provided http.Handler with JSON error and request logging.
//
// Parameters:
// - next: The http.Handler to be wrapped with logging.
//
// Returns:
// - http.Handler: The wrapped http.Handler.
func JSONLogger(next http.Handler) http.Handler {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, logHandlerOpts))
	accLogger := slog.New(slog.NewJSONHandler(os.Stdout, accLogHandlerOpts))

	return New(logger, accLogger)(next)
}

// GetLoggerFromContext extracts the logger from the provided context.
//
// Parameters:
// - ctx: The context from which to extract the logger.
//
// Returns:
// - *slog.Logger: The extracted logger.
// - bool: A boolean value indicating whether the logger was successfully extracted.
func GetLoggerFromContext(ctx context.Context) (*slog.Logger, bool) {
	l, ok := ctx.Value(loggerKey).(*slog.Logger)
	return l, ok
}

// GetLoggerFromRequest extracts the logger from the provided HTTP request's context.
//
// Parameters:
// - r: A pointer to the HTTP request from which to extract the logger.
//
// Returns:
// - *slog.Logger: The extracted logger.
// - bool: A boolean value indicating whether the logger was successfully extracted.
func GetLoggerFromRequest(r *http.Request) (*slog.Logger, bool) {
	ctx := r.Context()
	return GetLoggerFromContext(ctx)
}
