// Package persister is the port interface for local entry storage.
package persister

import (
	"context"

	"github.com/w-h-a/edge/internal/domain"
)

// Persister stores and retrieves entries on the local device.
type Persister interface {
	Put(ctx context.Context, entry domain.Entry) error
	Get(ctx context.Context, key string) (domain.Entry, error)
	Delete(ctx context.Context, key string) error
	List(ctx context.Context) ([]domain.Entry, error)
	ByScope(ctx context.Context, scope domain.Scope) ([]domain.Entry, error)
}
