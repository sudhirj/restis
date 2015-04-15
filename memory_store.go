package restis

import (
	"strconv"
)

type MemoryStore struct {
	strings map[string]string
	sets    map[string]map[string]bool
	hashes  map[string]map[string]string
	lists   map[string][]string
}

func (s *MemoryStore) Append(key, value string) int64 {
	s.strings[key] = s.strings[key] + value
	return int64(len(s.strings[key]))
}

func (s *MemoryStore) Get(key string) string {
	return s.strings[key]
}

func (s *MemoryStore) GetRange(key string, start, stop int64) string {
	start, stop = renormalize(int64(len(s.strings[key])), start, stop)
	return s.strings[key][start:stop]
}

func (s *MemoryStore) SetRange(key string, offset int64, value string) int64 {
	valueLength := int64(len(value))
	s.strings[key] = s.strings[key][:offset] + value + s.strings[key][offset+valueLength:]
	return s.Length(key)
}

func (s *MemoryStore) GetSet(key, value string) string {
	v := s.Get(key)
	s.Set(key, value)
	return v
}

func renormalize(length, start, stop int64) (int64, int64) {
	start = normalize(length, start)
	stop = normalize(length, stop)

	start = max(start, 0)
	stop = max(stop, -1)

	stop = stop + 1 // Make the stop index inclusive

	start = min(start, length)
	stop = min(stop, length)
	return start, stop
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
		if s.Exists(k) {
			m[k] = s.strings[k]
		}
	}
	return m
}

func (s *MemoryStore) MultiSet(data map[string]string) {
	for k, v := range data {
		s.strings[k] = v
	}
}

func (s *MemoryStore) MultiSetIfNotExists(data map[string]string) bool {
	for k, _ := range data {
		if s.Exists(k) {
			return false
		}
	}
	s.MultiSet(data)
	return true
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

func (s *MemoryStore) Length(key string) int64 {
	return int64(len(s.Get(key)))
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

func (s *MemoryStore) HashSetIfExists(key, field string, value string) bool {
	s.ensureHash(key)
	alreadyExists := s.HashExists(key, field)
	if alreadyExists {
		s.HashSet(key, field, value)
	}
	return alreadyExists
}

func (s *MemoryStore) HashSetIfNotExists(key, field string, value string) bool {
	s.ensureHash(key)
	alreadyExists := s.HashExists(key, field)
	if !alreadyExists {
		s.HashSet(key, field, value)
	}
	return !alreadyExists
}

func (s *MemoryStore) HashExists(key, field string) bool {
	s.ensureHash(key)
	_, exists := s.hashes[key][field]
	return exists
}

func (s *MemoryStore) HashMultiGet(key string, fields ...string) []string {
	values := []string{}
	for _, field := range fields {
		values = append(values, s.HashGet(key, field))
	}
	return values
}

func (s *MemoryStore) HashMultiSet(key string, data map[string]string) {
	for field, value := range data {
		s.HashSet(key, field, value)
	}
}

func (s *MemoryStore) HashLength(key string) int64 {
	return int64(len(s.hashes[key]))
}

func (s *MemoryStore) HashKeys(key string) []string {
	keys := []string{}
	for key, _ := range s.hashes[key] {
		keys = append(keys, key)
	}
	return keys
}

func (s *MemoryStore) HashValues(key string) []string {
	values := []string{}
	for _, value := range s.hashes[key] {
		values = append(values, value)
	}
	return values
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

func (s *MemoryStore) ListLeftPush(key string, values ...string) int64 {
	for _, value := range values {
		s.lists[key] = append([]string{value}, s.lists[key]...)
	}
	return s.ListLength(key)
}

func (s *MemoryStore) ListRightPush(key string, values ...string) int64 {
	for _, value := range values {
		s.lists[key] = append(s.lists[key], value)
	}
	return s.ListLength(key)
}

func (s *MemoryStore) ListLength(key string) int64 {
	return int64(len(s.lists[key]))
}

func (s *MemoryStore) ListLeftPop(key string) string {
	popped := s.lists[key][0]
	s.lists[key] = s.lists[key][1:]
	return popped
}

func (s *MemoryStore) ListRightPop(key string) string {
	lastIndex := s.ListLength(key) - 1
	popped := s.lists[key][lastIndex]
	s.lists[key] = s.lists[key][0:lastIndex]
	return popped
}

func min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func normalize(length, offset int64) int64 {
	if offset < 0 {
		return offset + length
	}
	return offset
}

func outOfBounds(length, index int64) bool {
	return index > length-1 || index < 0
}

func (s *MemoryStore) ListRange(key string, start, stop int64) []string {
	start, stop = renormalize(s.ListLength(key), start, stop)
	return s.lists[key][start:stop]
}

func (s *MemoryStore) ListSet(key string, index int64, value string) bool {
	length := s.ListLength(key)
	index = normalize(length, index)
	if outOfBounds(length, index) {
		return false
	}
	s.lists[key][index] = value
	return true
}

func (s *MemoryStore) ListIndex(key string, index int64) string {
	length := s.ListLength(key)
	index = normalize(length, index)
	if outOfBounds(length, index) {
		return ""
	}
	return s.lists[key][index]
}

func (s *MemoryStore) ListTrim(key string, start, stop int64) {
	s.lists[key] = s.ListRange(key, start, stop)
}

func NewMemoryStore() Store {
	return &MemoryStore{
		strings: make(map[string]string),
		sets:    make(map[string]map[string]bool),
		hashes:  make(map[string]map[string]string),
		lists:   make(map[string][]string),
	}
}
