package util

type Empty struct{}

// sets.String is a set of strings, implemented via map[string]struct{} for minimal memory consumption.
type StringSet map[string]Empty

// New creates a String from a list of values.
func NewStringSets(items ...string) StringSet {
	ss := StringSet{}
	ss.Insert(items...)
	return ss
}

// Insert adds items to the set.
func (s StringSet) Insert(items ...string) {
	for _, item := range items {
		s[item] = Empty{}
	}
}

// Delete removes all items from the set.
func (s StringSet) Delete(items ...string) {
	for _, item := range items {
		delete(s, item)
	}
}

// Has returns true if and only if item is contained in the set.
func (s StringSet) Has(item string) bool {
	_, contained := s[item]
	return contained
}

// HasAll returns true if and only if all items are contained in the set.
func (s StringSet) HasAll(items ...string) bool {
	for _, item := range items {
		if !s.Has(item) {
			return false
		}
	}
	return true
}

// HasAny returns true if any items are contained in the set.
func (s StringSet) HasAny(items ...string) bool {
	for _, item := range items {
		if s.Has(item) {
			return true
		}
	}
	return false
}
