package restis

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

type storeGenerator func() Store

func CheckStringOperations(t *testing.T, store StringStore) {
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

func CheckSetOperations(t *testing.T, store SetStore) {
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

func CheckHashOperations(t *testing.T, store HashStore) {
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

func CheckListOperations(t *testing.T, store ListStore) {
	assert.Equal(t, 1, store.ListRightPush("lk1", "lv1"))
	assert.Equal(t, 1, store.ListLength("lk1"))
	assert.Equal(t, 3, store.ListRightPush("lk1", "lv2", "lv3"))
	assert.Equal(t, 3, store.ListLength("lk1"))

	assert.Equal(t, []string{"lv1"}, store.ListRange("lk1", 0, 0))
	assert.Equal(t, []string{"lv1", "lv2"}, store.ListRange("lk1", 0, 1))
	assert.Equal(t, []string{"lv2", "lv3"}, store.ListRange("lk1", 1, 10))

	assert.Equal(t, []string{"lv1", "lv2", "lv3"}, store.ListRange("lk1", -3, 2))
	assert.Equal(t, []string{"lv1", "lv2", "lv3"}, store.ListRange("lk1", -100, 100))
	assert.Equal(t, []string{}, store.ListRange("lk1", 5, 10))

	assert.Equal(t, []string{"lv1", "lv2", "lv3"}, store.ListRange("lk1", 0, -1))
	assert.Equal(t, []string{"lv1", "lv2"}, store.ListRange("lk1", 0, -2))
	assert.Equal(t, []string{}, store.ListRange("lk1", 0, -20))
	assert.Equal(t, []string{}, store.ListRange("lk1", -10, -20))

	assert.Equal(t, 4, store.ListLeftPush("lk1", "lv0"))
	assert.Equal(t, 4, store.ListLength("lk1"))
	assert.Equal(t, []string{"lv0", "lv1", "lv2", "lv3"}, store.ListRange("lk1", 0, 10))
	assert.Equal(t, 6, store.ListLeftPush("lk1", "lv-1", "lv-2"))
	assert.Equal(t, []string{"lv-2", "lv-1", "lv0", "lv1", "lv2", "lv3"}, store.ListRange("lk1", 0, 10))

	assert.Equal(t, "lv-2", store.ListLeftPop("lk1"))
	assert.Equal(t, 5, store.ListLength("lk1"))
	assert.Equal(t, []string{"lv-1", "lv0", "lv1", "lv2", "lv3"}, store.ListRange("lk1", 0, 10))

	assert.Equal(t, "lv3", store.ListRightPop("lk1"))
	assert.Equal(t, 4, store.ListLength("lk1"))
	assert.Equal(t, []string{"lv-1", "lv0", "lv1", "lv2"}, store.ListRange("lk1", 0, 10))

	assert.False(t, store.ListSet("lk1", 34, "outofrange"))
	assert.Equal(t, []string{"lv-1", "lv0", "lv1", "lv2"}, store.ListRange("lk1", 0, 10))
	assert.True(t, store.ListSet("lk1", 0, "lv-1.2"))
	assert.Equal(t, []string{"lv-1.2", "lv0", "lv1", "lv2"}, store.ListRange("lk1", 0, 10))
	assert.False(t, store.ListSet("lk1", 4, "lv-1.2"))
	assert.Equal(t, []string{"lv-1.2", "lv0", "lv1", "lv2"}, store.ListRange("lk1", 0, 10))
	assert.True(t, store.ListSet("lk1", -1, "lv2.2"))
	assert.Equal(t, []string{"lv-1.2", "lv0", "lv1", "lv2.2"}, store.ListRange("lk1", 0, 10))

	assert.False(t, store.ListSet("lk1", -5, "oob"))
	assert.Equal(t, []string{"lv-1.2", "lv0", "lv1", "lv2.2"}, store.ListRange("lk1", 0, 10))
	assert.Equal(t, "lv-1.2", store.ListIndex("lk1", 0))
	assert.Equal(t, "lv0", store.ListIndex("lk1", 1))
	assert.Equal(t, "lv2.2", store.ListIndex("lk1", -1))
	assert.Equal(t, "", store.ListIndex("lk1", -10))

	store.ListTrim("lk1", 1, -2)
	assert.Equal(t, []string{"lv0", "lv1"}, store.ListRange("lk1", 0, 10))

}

func RunAllTestsOnStore(t *testing.T, storeGen storeGenerator) {
	CheckStringOperations(t, storeGen())
	CheckSetOperations(t, storeGen())
	CheckHashOperations(t, storeGen())
	CheckListOperations(t, storeGen())
}
