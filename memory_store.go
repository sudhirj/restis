package main

import "strconv"

// MemoryStore acts as a datastore using the current instance's memory.
// Does not offer persistence or distribution. Works only in single instance setups.
type MemoryStore struct {
	strings map[string]string
	sets    map[string]map[string]bool
	hashes  map[string]map[string]string
}

func (s *MemoryStore) Get(key string) string {
	return s.strings[key]
}

func (s *MemoryStore) Set(key string, value string) {
	s.strings[key] = value
}

func (s *MemoryStore) SetIfNotExists(key string, value string) bool {
	alreadyExists := s.Exists(key)
	if !alreadyExists {
		s.Set(key, value)
	}
	return !alreadyExists
}

func (s *MemoryStore) SetIfExists(key string, value string) bool {
	alreadyExists := s.Exists(key)
	if alreadyExists {
		s.Set(key, value)
	}
	return alreadyExists
}

func (s *MemoryStore) MultiGet(keys []string) map[string]string {
	m := map[string]string{}
	for _, k := range keys {
		m[k] = s.strings[k]
	}
	return m
}

func (s *MemoryStore) MultiSet(data map[string]string) {
	for k, v := range data {
		s.strings[k] = v
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
	_, exists := s.strings[key]
	return exists
}

func (s *MemoryStore) ensureSet(key string) {
	if _, ok := s.sets[key]; !ok {
		s.sets[key] = make(map[string]bool)
	}
}

func (s *MemoryStore) SetAdd(key string, values ...string) {
	s.ensureSet(key)
	for _, value := range values {
		s.sets[key][value] = true
	}
}

func (s *MemoryStore) SetRemove(key string, values ...string) {
	s.ensureSet(key)
	for _, value := range values {
		delete(s.sets[key], value)
	}
}

func (s *MemoryStore) SetIsMember(key string, value string) bool {
	_, exists := s.sets[key][value]
	return exists
}

func (s *MemoryStore) SetCardinality(key string) int64 {
	s.ensureSet(key)
	return int64(len(s.sets[key]))
}

func (s *MemoryStore) SetMembers(key string) []string {
	s.ensureSet(key)
	values := []string{}
	for val := range s.sets[key] {
		values = append(values, val)
	}
	return values
}

func (s *MemoryStore) ensureHash(key string) {
	if _, ok := s.hashes[key]; !ok {
		s.hashes[key] = make(map[string]string)
	}
}
func (s *MemoryStore) HashGet(key, field string) string {
	s.ensureHash(key)
	return s.hashes[key][field]
}
func (s *MemoryStore) HashSet(key, field, value string) {
	s.ensureHash(key)
	s.hashes[key][field] = value
}

func (s *MemoryStore) transformNumber(key string, transform func(int64) int64) int64 {
	n, err := strconv.ParseInt(s.strings[key], 10, 64)
	if err != nil {
		n = 0
	}
	n = transform(n)
	s.strings[key] = strconv.FormatInt(n, 10)
	return n
}

// NewMemoryStore creates a new memory store with a string map
func NewMemoryStore() Store {
	return &MemoryStore{
		strings: make(map[string]string),
		sets:    make(map[string]map[string]bool),
		hashes:  make(map[string]map[string]string),
	}
}
