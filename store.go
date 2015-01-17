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
	AddToSet(key string, values ...string)
	RemoveFromSet(key string, values ...string)
	IsMemberOfSet(key string, value string) bool
	CardinalityOfSet(key string) int64
	MembersOfSet(key string) []string
}
