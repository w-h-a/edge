// Package service contains the controller layer.
package service

// SyncManager runs background sync. Periodic, retry with backoff.
// Connects to delta via transport port, exchanges deltas, confirms convergence.
type SyncManager struct {
	// ports injected at construction
}
