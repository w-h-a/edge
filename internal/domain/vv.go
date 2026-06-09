package domain

import "maps"

// VersionVector is a PLACEHOLDER. It will be replaced by an import of
// github.com/w-h-a/meld/crdt/versionvector. The meld type
// uses a sorted-slice representation for cache-friendly O(n+m) merge.
//
// Once the meld library is used, this file should be deleted and
// edge.Entry.Version should hold a serialized meld version vector
// via its Marshal. Anywhere this VersionVector is used in edge
// should switch to meld VersionVector.
//
// Do not add methods or extend this type. Any feature added here will need
// to be removed during the meld swap.
type VersionVector struct {
	Entries map[string]uint64 // nodeID -> counter
}

// Increment advances this node's counter. Placeholder. Replace with
// meld VersionVector.
func (vv VersionVector) Increment(nodeID string) VersionVector {
	next := make(map[string]uint64, len(vv.Entries))
	maps.Copy(next, vv.Entries)
	next[nodeID] = next[nodeID] + 1
	return VersionVector{Entries: next}
}

// Merge returns the pointwise max of two version vectors. Placeholder.
// Replace with meld VersionVector.
func (vv VersionVector) Merge(other VersionVector) VersionVector {
	merged := make(map[string]uint64, len(vv.Entries))
	maps.Copy(merged, vv.Entries)
	for k, v := range other.Entries {
		if v > merged[k] {
			merged[k] = v
		}
	}
	return VersionVector{Entries: merged}
}
