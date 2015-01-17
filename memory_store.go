package main

import "strconv"

// MemoryStore acts as a datastore using the current instance's memory.
// Does not offer persistence or distribution. Works only in single instance setups.
type MemoryStore struct {
	keyValues map[string]string
	sets      map[string]map[string]bool
}

func (s *MemoryStore) Get(key string) string {
	return s.keyValues[key]
}

func (s *MemoryStore) Set(key string, value string) {
	s.keyValues[key] = value
}

func (s *MemoryStore) SetNX(key string, value string) bool {
	alreadyExists := s.Exists(key)
	if !alreadyExists {
		s.Set(key, value)
	}
	return !alreadyExists
}

func (s *MemoryStore) SetEX(key string, value string) bool {
	alreadyExists := s.Exists(key)
	if alreadyExists {
		s.Set(key, value)
	}
	return alreadyExists
}

func (s *MemoryStore) MultiGet(keys []string) map[string]string {
	m := map[string]string{}
	for _, k := range keys {
		m[k] = s.keyValues[k]
	}
	return m
}

func (s *MemoryStore) MultiSet(data map[string]string) {
	for k, v := range data {
		s.keyValues[k] = v
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
	_, exists := s.keyValues[key]
	return exists
}

func (s *MemoryStore) ensure(key string) {
	if _, ok := s.sets[key]; !ok {
		s.sets[key] = make(map[string]bool)
	}
}

func (s *MemoryStore) SAdd(key string, values ...string) {
	s.ensure(key)
	for _, value := range values {
		s.sets[key][value] = true
	}
}

func (s *MemoryStore) SIsMember(key string, value string) bool {
	_, exists := s.sets[key][value]
	return exists
}

func (s *MemoryStore) transformNumber(key string, transform func(int64) int64) int64 {
	n, err := strconv.ParseInt(s.keyValues[key], 10, 64)
	if err != nil {
		n = 0
	}
	n = transform(n)
	s.keyValues[key] = strconv.FormatInt(n, 10)
	return n
}

// NewMemoryStore creates a new memory store with a string map
func NewMemoryStore() Store {
	return &MemoryStore{
		keyValues: make(map[string]string),
		sets:      make(map[string]map[string]bool),
	}
}
