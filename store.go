package restis

// Strings is the list of operations available for string keys
type StringStore interface {
	Append(key, value string) int64
	Get(key string) string
	GetRange(key string, start, stop int64) string
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
}

// Sets is the list of operations available for set keys
type SetStore interface {
	SetAdd(key string, values ...string)
	SetRemove(key string, values ...string)
	SetIsMember(key string, value string) bool
	SetMembers(key string) []string
	SetCardinality(key string) int64
}

// Hashes is the list of operations available for hash types
type HashStore interface {
	HashGet(key, field string) string
	HashSet(key, field, value string)
	HashLength(key string) int64
	HashMultiGet(key string, fields ...string) []string
	HashMultiSet(key string, data map[string]string)
	HashExists(key, field string) bool
	HashKeys(key string) []string
	HashValues(key string) []string
	HashSetIfExists(key, field string, value string) bool
	HashSetIfNotExists(key, field string, value string) bool
}

type ListStore interface {
	ListLeftPush(key string, values ...string) int64
	ListRightPush(key string, values ...string) int64
	ListLeftPop(key string) string
	ListRightPop(key string) string
	ListLength(key string) int64
	ListRange(key string, start, stop int64) []string
	ListSet(key string, index int64, value string) bool
	ListIndex(key string, index int64) string
	ListTrim(key string, start, stop int64)
}

// Store is full interface for the any datastore that backs Restis
type Store interface {
	StringStore
	SetStore
	HashStore
	ListStore
}
