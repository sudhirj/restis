package main

import "testing"
import "github.com/stretchr/testify/assert"
import "sort"

func StringOperations(t *testing.T, store Store) {
	store.Set("k1", "v1")
	assert.Equal(t, "v1", store.Get("k1"))

	store.Set("k2", "v2")
	assert.Equal(t, "v1", store.Get("k1"))
	assert.Equal(t, "v2", store.Get("k2"))

	store.Set("k1", "v1.1")
	assert.Equal(t, "v1.1", store.Get("k1"))

	assert.Equal(t, map[string]string{"k1": "v1.1", "k2": "v2"}, store.MultiGet([]string{"k1", "k2"}))
	store.MultiSet(map[string]string{"k1": "v1.2", "k2": "v2.1"})
	assert.Equal(t, map[string]string{"k1": "v1.2", "k2": "v2.1"}, store.MultiGet([]string{"k1", "k2"}))

	assert.Equal(t, 1, store.Increment("n1"))
	assert.Equal(t, "1", store.Get("n1"))
	assert.Equal(t, 2, store.Increment("n1"))
	assert.Equal(t, 3, store.Increment("n1"))
	assert.Equal(t, 2, store.Decrement("n1"))
	assert.Equal(t, 1, store.Decrement("n1"))
	assert.Equal(t, "1", store.Get("n1"))
	assert.Equal(t, 5, store.IncrementBy("n1", 4))
	assert.Equal(t, 3, store.DecrementBy("n1", 2))

	assert.True(t, store.Exists("n1"))
	assert.False(t, store.Exists("non existent key"))

	assert.Equal(t, "", store.Get("non existent key"))

	assert.False(t, store.SetIfExists("ek1", "ev1"))
	assert.Equal(t, "", store.Get("ek1"))
	store.Set("ek1", "some old value")
	assert.True(t, store.SetIfExists("ek1", "ev2"))
	assert.Equal(t, "ev2", store.Get("ek1"))

	assert.True(t, store.SetIfNotExists("nk1", "vx1"))
	assert.Equal(t, "vx1", store.Get("nk1"))
	assert.False(t, store.SetIfNotExists("nk1", "vx2"))
	assert.Equal(t, "vx1", store.Get("nk1"))
}

func SetOperations(t *testing.T, store Store) {
	assert.False(t, store.SetIsMember("sk1", "v1"))
	assert.Equal(t, 0, store.SetCardinality("sk1"))

	store.SetAdd("sk1", "v1")
	assert.True(t, store.SetIsMember("sk1", "v1"))
	assert.False(t, store.SetIsMember("sk1", "v2"))
	assert.Equal(t, 1, store.SetCardinality("sk1"))

	store.SetAdd("sk1", "v2", "v1")
	assert.True(t, store.SetIsMember("sk1", "v2"))
	assert.Equal(t, 2, store.SetCardinality("sk1"))

	members := store.SetMembers("sk1")
	sort.Sort(sort.StringSlice(members))

	expected := []string{"v1", "v2"}
	sort.Sort(sort.StringSlice(expected))
	assert.Equal(t, members, expected)

	store.SetRemove("sk1", "v1")
	assert.False(t, store.SetIsMember("sk1", "v1"))
	assert.Equal(t, 1, store.SetCardinality("sk1"))
}

func HashOperations(t *testing.T, store Store) {
	store.HashSet("hk1", "f1", "v1")
	assert.Equal(t, "v1", store.HashGet("hk1", "f1"))
	assert.Equal(t, "", store.HashGet("hk1", "f2"))
}

func TestMemoryStore(t *testing.T) {
	memoryStore := NewMemoryStore()
	StringOperations(t, memoryStore)
	SetOperations(t, memoryStore)
	HashOperations(t, memoryStore)
}
