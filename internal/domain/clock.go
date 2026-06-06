package domain

// VersionVector is a PLACEHOLDER. It will be replaced by an import of
// github.com/w-h-a/meld/crdt/vclock when meld-tuwl lands. The meld type
// uses a sorted-slice representation (Syncthing pattern) for cache-friendly
// O(n+m) merge. The map-backed implementation here exists only so edge
// compiles before meld ships.
//
// When meld-tuwl ships, this file should be deleted and edge.Entry.Version
// should hold a serialized meld vclock via its MarshalBinary/UnmarshalBinary.
// Anywhere this VersionVector is used in edge (DeltaCompute, session.go,
// transport.go) should switch to vclock.VersionVector.
//
// Do not add methods or extend this type. Any feature added here will need
// to be removed during the meld swap. See meld-tuwl for the planned shape.
type VersionVector struct {
	Entries map[string]uint64 // nodeID -> counter
}

// Increment advances this node's counter. Placeholder. Replace with
// meld vclock.VersionVector.Increment when meld-tuwl lands.
func (vv VersionVector) Increment(nodeID string) VersionVector {
	next := make(map[string]uint64, len(vv.Entries))
	for k, v := range vv.Entries {
		next[k] = v
	}
	next[nodeID] = next[nodeID] + 1
	return VersionVector{Entries: next}
}

// Merge returns the pointwise max of two version vectors. Placeholder.
// Replace with meld vclock.VersionVector.Merge when meld-tuwl lands.
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
