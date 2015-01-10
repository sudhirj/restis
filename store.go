package main

// Store is full interface for the any datastore that backs Restis
type Store interface {
	Get(key string) string
	Set(key string, value string)
	MultiGet(keys []string) map[string]string
	MultiSet(map[string]string)
	Increment(key string) int64
	Decrement(key string) int64
	IncrementBy(key string, delta int64) int64
	DecrementBy(key string, delta int64) int64
	Exists(key string) bool
}