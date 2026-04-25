//go:build !nocaddy

package nats

import (
	"context"

	"github.com/caddyserver/caddy/v2"
	"go.uber.org/zap"
)

// getContext returns the active Caddy context so request-scoped operations
// (e.g. nats.RequestWithContext) honour Caddy's lifecycle. Falls back to
// context.Background() when no Caddy server is running (e.g. php-cli),
// where caddy.ActiveContext() returns a zero Context with a nil embedded
// context.Context that would panic on Deadline()/Done() calls.
func getContext() context.Context {
	ctx := caddy.ActiveContext()
	if ctx.Context == nil {
		return context.Background()
	}
	return ctx
}

// getLogger returns the Caddy-managed logger for this package.
func getLogger() *zap.Logger {
	return caddy.Log().Named("nats")
}
