package domain

// Entry is a context record with CRDT metadata.
// Different entry types use different merge strategies:
// LWW-Register for preferences, OR-Set for accumulated patterns.
type Entry struct {
	Key      string
	Value    []byte
	Scope    Scope
	Strategy MergeStrategy
	Version  []byte // serialized version vector from meld
}

// MergeStrategy tags which CRDT merge to use for this entry.
type MergeStrategy int

const (
	LWW MergeStrategy = iota // last-writer-wins (preferences, corrections)
	Set                      // OR-Set (accumulated patterns)
)

// Merge combines two entries using the tagged CRDT strategy.
// Commutative, associative, idempotent.
func Merge(a, b Entry) Entry {
	return Entry{} // TODO: dispatch on Strategy, delegate to meld
}
