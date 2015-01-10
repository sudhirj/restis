package main

import "strconv"

// MemoryStore acts as a datastore using the current instance's memory.
// Does not offer persistence or distribution. Works only in single instance setups.
type MemoryStore struct {
	stringMap map[string]string
}

func (s *MemoryStore) Get(key string) string {
	return s.stringMap[key]
}

func (s *MemoryStore) Set(key string, value string) {
	s.stringMap[key] = value
}

func (s *MemoryStore) MultiGet(keys []string) map[string]string {
	m := map[string]string{}
	for _, k := range keys {
		m[k] = s.stringMap[k]
	}
	return m
}

func (s *MemoryStore) MultiSet(data map[string]string) {
	for k, v := range data {
		s.stringMap[k] = v
	}
}

func (s *MemoryStore) Increment(key string) int64 {
	return s.transformNumber(key, func(n int64) int64 { return n + 1 })
}

func (s *MemoryStore) Decrement(key string) int64 {
	return s.transformNumber(key, func(n int64) int64 { return n - 1 })
}

func (s *MemoryStore) IncrementBy(key string, delta int64) int64 {
	return s.transformNumber(key, func(n int64) int64 { return n + delta })
}

func (s *MemoryStore) DecrementBy(key string, delta int64) int64 {
	return s.transformNumber(key, func(n int64) int64 { return n - delta })
}

func (s *MemoryStore) Exists(key string) bool {
	_, exists := s.stringMap[key]
	return exists
}

func (s *MemoryStore) transformNumber(key string, transform func(int64) int64) int64 {
	n, err := strconv.ParseInt(s.stringMap[key], 10, 64)
	if err != nil {
		n = 0
	}
	n = transform(n)
	s.stringMap[key] = strconv.FormatInt(n, 10)
	return n
}

// NewMemoryStore creates a new memory store with a string map
func NewMemoryStore() Store {
	return &MemoryStore{
		stringMap: map[string]string{},
	}
}
