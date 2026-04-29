// Package transport is the port interface for sync communication with delta.
package transport

import (
	"context"

	"github.com/w-h-a/edge/internal/domain"
)

// Transport connects to a delta node and exchanges deltas.
type Transport interface {
	Connect(ctx context.Context, addr string) error
	Negotiate(ctx context.Context, localVV domain.VersionVector) (domain.VersionVector, error)
	SendDelta(ctx context.Context, entries []domain.Entry) error
	ReceiveDelta(ctx context.Context) ([]domain.Entry, error)
	Close(ctx context.Context) error
}
