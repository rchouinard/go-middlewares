package middlewares

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	loggerKey ctxKey = "logger"
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

// Logger implements a middleware which adds a text [slog] access and error logger to each request.
//
// The access logger outputs to Stdout while the error logger outputs to Stderr.
func Logger(next http.Handler) http.Handler {
	logger := slog.New(slog.NewTextHandler(os.Stderr, logHandlerOpts))
	accLogger := slog.New(slog.NewTextHandler(os.Stdout, accLogHandlerOpts))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logRequest(w, r, next, logger, accLogger)
	})
}

// JSONLogger implements a middleware which adds a JSON [slog] access and error logger to each request.
//
// The access logger outputs to Stdout while the error logger outputs to Stderr.
func JSONLogger(next http.Handler) http.Handler {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, logHandlerOpts))
	accLogger := slog.New(slog.NewJSONHandler(os.Stdout, accLogHandlerOpts))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logRequest(w, r, next, logger, accLogger)
	})
}

// NewLogger attaches the passed in [slog.Logger] access and error loggers to a request.
func NewLogger(logger, accLogger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logRequest(w, r, next, logger, accLogger)
		})
	}
}

// GetLoggerFromContext returns the error logger from a given request context.
func GetLoggerFromContext(ctx context.Context) (*slog.Logger, bool) {
	l, ok := ctx.Value(loggerKey).(*slog.Logger)
	return l, ok
}

// GetLoggerFromRequest returns the error logger from a given request.
func GetLoggerFromRequest(r *http.Request) (*slog.Logger, bool) {
	ctx := r.Context()
	return GetLoggerFromContext(ctx)
}

// logRequest parses the request and response data to write out an access log line.
func logRequest(w http.ResponseWriter, r *http.Request, next http.Handler, logger, accLogger *slog.Logger) {
	// create new request context with logger attached
	// this will be passed down to the remaining handlers
	ctx := context.WithValue(r.Context(), loggerKey, logger)

	// run the remaining handlers and gather stats
	lrw := NewResponseWriter(w)
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
	if id, ok := GetRequestIDFromRequest(r); ok {
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
}
