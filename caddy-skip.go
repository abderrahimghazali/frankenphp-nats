//go:build nocaddy

package nats

import (
	"context"

	"go.uber.org/zap"
)

// getContext returns a background context when Caddy is excluded from the
// build (-tags nocaddy). Suitable for unit tests or non-Caddy embedders.
func getContext() context.Context {
	return context.Background()
}

// getLogger returns a no-op logger when Caddy is excluded from the build.
func getLogger() *zap.Logger {
	return zap.NewNop()
}
