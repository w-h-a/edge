package domain

// DeltaCompute returns the entries that the remote side needs.
// Given my version vector, their version vector, and my entries,
// return entries where my version exceeds theirs. Pure function.
func DeltaCompute(myVV, theirVV VersionVector, myEntries []Entry) []Entry {
	return nil // TODO: compare vectors, filter entries
}
