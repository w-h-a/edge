package domain

// Scope namespaces entries. One edge instance per device, all context
// in one store. Repo is metadata, not a partition boundary.
type Scope struct {
	Global  bool   // portfolio-wide (e.g. ideation context)
	Project string // project-scoped (e.g. tally decisions)
	Repo    string // repo-scoped (e.g. implementation tasks)
}

// Contains returns whether an entry belongs to this scope. Pure predicate.
func (s Scope) Contains(entry Entry) bool {
	return false // TODO
}
