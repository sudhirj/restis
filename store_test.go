package main

import "testing"
import "github.com/stretchr/testify/assert"

func StringOperations(t *testing.T, store Store) {
	store.Set("k1", "v1")
	assert.Equal(t, "v1", store.Get("k1"), "Get and set didn't work.")

	store.Set("k2", "v2")
	assert.Equal(t, "v1", store.Get("k1"), "Set overrode a wrong variable.")
	assert.Equal(t, "v2", store.Get("k2"), "Set and Get didn't work.")

	store.Set("k1", "v1.1")
	assert.Equal(t, "v1.1", store.Get("k1"), "Set didn't override.")

	assert.Equal(t, map[string]string{"k1": "v1.1", "k2": "v2"}, store.MultiGet([]string{"k1", "k2"}), "MultiGet doesn't work")
	store.MultiSet(map[string]string{"k1": "v1.2", "k2": "v2.1"})
	assert.Equal(t, map[string]string{"k1": "v1.2", "k2": "v2.1"}, store.MultiGet([]string{"k1", "k2"}), "MultiGet doesn't work")

	assert.Equal(t, 1, store.Increment("n1"), "should be 1 for non existent keys")
	assert.Equal(t, "1", store.Get("n1"), "should have saved number in key")
	assert.Equal(t, 2, store.Increment("n1"), "should have saved increment")
	assert.Equal(t, 3, store.Increment("n1"), "should have saved increment")
	assert.Equal(t, 2, store.Decrement("n1"), "should decrement")
	assert.Equal(t, 1, store.Decrement("n1"), "should decrement")
	assert.Equal(t, "1", store.Get("n1"), "should have saved number in key")
}

func TestMemoryStore(t *testing.T) {
	memoryStore := NewMemoryStore()
	StringOperations(t, memoryStore)
}
