package restis

import "testing"
import "github.com/stretchr/testify/assert"
import "sort"

func StringOperations(t *testing.T, store StringStore) {
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

func SetOperations(t *testing.T, store SetStore) {
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

func HashOperations(t *testing.T, store HashStore) {
	store.HashSet("hk1", "f1", "v1")
	assert.Equal(t, "v1", store.HashGet("hk1", "f1"))
	assert.Equal(t, "", store.HashGet("hk1", "f2"))
	assert.False(t, store.HashExists("hk1", "f0"))
	assert.True(t, store.HashExists("hk1", "f1"))
	assert.False(t, store.HashExists("hk2", "f1"))

	store.HashSet("hm1", "k1", "v1")
	store.HashSet("hm1", "k2", "v2")
	assert.Equal(t, []string{"v1", "v2", ""}, store.HashMultiGet("hm1", "k1", "k2", "unknownkey"))
	assert.Equal(t, []string{"v1", "v2", ""}, store.HashMultiGet("hm1", []string{"k1", "k2", "unknownkey"}...))
	store.HashMultiSet("hm2", map[string]string{"k1": "v1", "k2": "v2"})
	assert.Equal(t, "v1", store.HashGet("hm2", "k1"))
	assert.Equal(t, "v2", store.HashGet("hm2", "k2"))
	assert.Equal(t, 2, store.HashLength("hm2"))

	assert.Equal(t, 0, store.HashLength("nonexistenthash"))
	keys := store.HashKeys("hm2")
	sort.Strings(keys)
	assert.Equal(t, []string{"k1", "k2"}, keys)
	values := store.HashValues("hm2")
	sort.Strings(values)
	assert.Equal(t, []string{"v1", "v2"}, values)
	assert.Equal(t, 2, store.HashLength("hm2"))

	assert.False(t, store.HashSetIfExists("hm2", "k3", "v3"))
	store.HashSet("hm2", "k3", "v3")
	assert.True(t, store.HashSetIfExists("hm2", "k3", "v3.1"))
	assert.Equal(t, "v3.1", store.HashGet("hm2", "k3"))

	assert.True(t, store.HashSetIfNotExists("hm2", "k4", "v4"))
	assert.Equal(t, "v4", store.HashGet("hm2", "k4"))
	assert.False(t, store.HashSetIfNotExists("hm2", "k4", "v4.1"))
	assert.Equal(t, "v4", store.HashGet("hm2", "k4"))
	assert.Equal(t, 4, store.HashLength("hm2"))

}

func ListOperations(t *testing.T, store ListStore) {
	assert.Equal(t, 1, store.ListPush("lk1", "lv1"))
	assert.Equal(t, 1, store.ListLength("lk1"))
}

func RunAllTestsOnStore(t *testing.T, store Store) {
	StringOperations(t, store)
	SetOperations(t, store)
	HashOperations(t, store)
	ListOperations(t, store)
}

func TestMemoryStore(t *testing.T) {
	memoryStore := NewMemoryStore()
	RunAllTestsOnStore(t, memoryStore)
}
