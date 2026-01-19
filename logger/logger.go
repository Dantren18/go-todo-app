package logger

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type TraceIDKey struct{}

func ContextLogger(ctx context.Context) *slog.Logger {
	traceID := ctx.Value(TraceIDKey{})
	if traceID == nil {
		return slog.Default()
	}
	return slog.With("trace_id", traceID.(string))
}

// Middleware to inject trace ID
func TraceIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := fmt.Sprintf("%d", time.Now().UnixNano())
		ctx := context.WithValue(r.Context(), TraceIDKey{}, traceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
