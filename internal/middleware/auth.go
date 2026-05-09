package middleware

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"qc/config"
)

type contextKey string

const ExternalIpKey contextKey = "external_ip"
const AuthErrorKey contextKey = "auth_error"

func AuthRequired(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if err := r.ParseForm(); err != nil {
				slog.Warn(
					"vote access check failed",
					slog.String("reason", "invalid_form"),
					slog.String("remote_addr", r.RemoteAddr),
					slog.Any("error", err),
				)
				ctx := context.WithValue(r.Context(), AuthErrorKey, "invalid_form")
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			externalIp := extractExternalIp(r)
			slog.Info(
				"vote access attempt",
				slog.String("external_ip", externalIp),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
			)

			slog.Info(
				"vote access allowed",
				slog.String("external_ip", externalIp),
				slog.String("path", r.URL.Path),
			)

			ctx := context.WithValue(r.Context(), ExternalIpKey, extractExternalIp(r))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractExternalIp(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")

	if ip == "" {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			host = r.RemoteAddr
		}
		ip = host
	}

	return ip
}
