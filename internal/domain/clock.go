package domain

// VersionVector wraps meld/crdt/vclock for causal ordering.
// No wall clocks. Version vectors track which writes each replica has seen.
type VersionVector struct {
	Entries map[string]uint64 // nodeID -> counter
}

// Increment advances this node's counter.
func (vv VersionVector) Increment(nodeID string) VersionVector {
	next := make(map[string]uint64, len(vv.Entries))
	for k, v := range vv.Entries {
		next[k] = v
	}
	next[nodeID] = next[nodeID] + 1
	return VersionVector{Entries: next}
}

// Merge returns the pointwise max of two version vectors.
func (vv VersionVector) Merge(other VersionVector) VersionVector {
	merged := make(map[string]uint64, len(vv.Entries))
	for k, v := range vv.Entries {
		merged[k] = v
	}
	for k, v := range other.Entries {
		if v > merged[k] {
			merged[k] = v
		}
	}
	return VersionVector{Entries: merged}
}
