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
	loggerKey string = "logger"
)

type LoggerConfig struct {
	AccessHandler slog.Handler
	ErrorHandler  slog.Handler
}

var defaultLoggerConfig = LoggerConfig{
	AccessHandler: slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.MessageKey || a.Key == slog.LevelKey {
				return slog.Attr{} // drop attribute
			}

			return a
		},
	}),
	ErrorHandler: slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelError,
	}),
}

func LoggerWithConfig(cfg LoggerConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var attrs []slog.Attr
			if id, ok := RequestIDFromRequest(r); ok {
				attrs = append(attrs, slog.String("request_id", id))
			}

			// add error logger to context
			errLog := slog.New(cfg.ErrorHandler.WithAttrs(attrs))
			ctx := context.WithValue(r.Context(), loggerKey, errLog)

			// run the remaining handlers and gather stats
			lrw := NewResponseWriter(w)
			start := time.Now()
			next.ServeHTTP(lrw, r.WithContext(ctx))
			duration := time.Since(start)

			addr, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				addr = r.RemoteAddr
			}

			accLog := slog.New(cfg.AccessHandler)
			f := append(attrs,
				slog.String("remote_ip", addr),
				slog.String("request", fmt.Sprintf("%s %s %s", r.Method, r.URL.String(), r.Proto)),
				slog.Int("response", lrw.Status()),
				slog.Int("bytes", lrw.Size()),
				slog.String("referer", r.Referer()),
				slog.String("agent", r.UserAgent()),
				slog.Duration("duration", duration),
			)

			accLog.LogAttrs(context.Background(), slog.LevelInfo, "access_log", f...)
		})
	}
}

func Logger(next http.Handler) http.HandlerFunc {
	return LoggerWithConfig(defaultLoggerConfig)(next).(http.HandlerFunc)
}

func LoggerFromContext(ctx context.Context) (*slog.Logger, bool) {
	l, ok := ctx.Value(loggerKey).(*slog.Logger)
	return l, ok
}

func LoggerFromRequest(r *http.Request) (*slog.Logger, bool) {
	ctx := r.Context()
	return LoggerFromContext(ctx)
}
