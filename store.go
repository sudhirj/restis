package main

// Store is full interface for the any datastore that backs Restis
type Store interface {
	// Strings
	Get(key string) string
	Set(key string, value string)
	SetIfExists(key string, value string) bool
	SetIfNotExists(key string, value string) bool
	MultiGet(keys []string) map[string]string
	MultiSet(map[string]string)
	Increment(key string) int64
	Decrement(key string) int64
	IncrementBy(key string, delta int64) int64
	DecrementBy(key string, delta int64) int64
	Exists(key string) bool

	// Sets
	SetAdd(key string, values ...string)
	SetRemove(key string, values ...string)
	SetIsMember(key string, value string) bool
	SetCardinality(key string) int64
	SetMembers(key string) []string

	// Hashes
	HashGet(key, field string) string
	HashSet(key, field, value string)
	HashExists(key, field string) bool
}
