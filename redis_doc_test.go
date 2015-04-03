package restis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func RunAllRedisDocChecksOnStore(t *testing.T, storeGen storeGenerator) {
	APPEND(t, storeGen())
	DECR(t, storeGen())
	DECRBY(t, storeGen())
	GET(t, storeGen())
	GETRANGE(t, storeGen())
	GETSET(t, storeGen())
	INCR(t, storeGen())
	INCRBY(t, storeGen())
	MGET(t, storeGen())
	MSET(t, storeGen())
	MSETNX(t, storeGen())
	SET(t, storeGen())
	SETNX(t, storeGen())
	SETRANGE(t, storeGen())
	STRLEN(t, storeGen())
}

func APPEND(t *testing.T, store StringStore) {
	assert.False(t, store.Exists("mykey"))
	assert.Equal(t, 5, store.Append("mykey", "Hello"))
	assert.Equal(t, 11, store.Append("mykey", " World"))
	assert.Equal(t, "Hello World", store.Get("mykey"))
}

func DECR(t *testing.T, store StringStore) {
	store.Set("mykey", "10")
	assert.Equal(t, 9, store.Decrement("mykey"))
	store.Set("mykey", "234293482390480948029348230948")
	assert.Equal(t, -1, store.Decrement("mykey")) // DEVIATION: Redis throws an error instead.
}

func DECRBY(t *testing.T, store StringStore) {
	store.Set("mykey", "10")
	assert.Equal(t, 7, store.DecrementBy("mykey", 3))
}

func GET(t *testing.T, store StringStore) {
	assert.Equal(t, "", store.Get("nonexisting"))
	store.Set("mykey", "Hello")
	assert.Equal(t, "Hello", store.Get("mykey"))
}

func GETRANGE(t *testing.T, store StringStore) {
	store.Set("mykey", "This is a string")
	assert.Equal(t, "This", store.GetRange("mykey", 0, 3))
	assert.Equal(t, "ing", store.GetRange("mykey", -3, -1))
	assert.Equal(t, "This is a string", store.GetRange("mykey", 0, -1))
	assert.Equal(t, "string", store.GetRange("mykey", 10, 100))
}

func GETSET(t *testing.T, store StringStore) {
	assert.Equal(t, 1, store.Increment("mycounter"))
	assert.Equal(t, "1", store.GetSet("mycounter", "0"))
	assert.Equal(t, "0", store.Get("mycounter"))

	store.Set("mykey", "Hello")
	assert.Equal(t, "Hello", store.GetSet("mykey", "World"))
	assert.Equal(t, "World", store.Get("mykey"))
}

func INCR(t *testing.T, store StringStore) {
	store.Set("mykey", "10")
	assert.Equal(t, 11, store.Increment("mykey"))
	assert.Equal(t, "11", store.Get("mykey"))
}

func INCRBY(t *testing.T, store StringStore) {
	store.Set("mykey", "10")
	assert.Equal(t, 15, store.IncrementBy("mykey", 5))
}

func MGET(t *testing.T, store StringStore) {
	store.Set("key1", "Hello")
	store.Set("key2", "World")
	assert.Equal(t, map[string]string{"key1": "Hello", "key2": "World"}, store.MultiGet([]string{"key1", "key2", "nonexisting"}))
}

func MSET(t *testing.T, store StringStore) {
	store.MultiSet(map[string]string{"key1": "Hello", "key2": "World"})
	assert.Equal(t, "Hello", store.Get("key1"))
	assert.Equal(t, "World", store.Get("key2"))
}

func MSETNX(t *testing.T, store StringStore) {
	assert.True(t, store.MultiSetIfNotExists(map[string]string{"key1": "Hello", "key2": "there"}))
	assert.False(t, store.MultiSetIfNotExists(map[string]string{"key2": "there", "key3": "world"}))
	assert.Equal(t, map[string]string{"key1": "Hello", "key2": "there"}, store.MultiGet([]string{"key1", "key2", "key3"}))
}

func SET(t *testing.T, store StringStore) {
	store.Set("mykey", "Hello")
	assert.Equal(t, "Hello", store.Get("mykey"))
}

func SETNX(t *testing.T, store StringStore) {
	assert.True(t, store.SetIfNotExists("mykey", "Hello"))
	assert.False(t, store.SetIfNotExists("mykey", "World"))
	assert.Equal(t, "Hello", store.Get("mykey"))
}

func SETRANGE(t *testing.T, store StringStore) {
}

func STRLEN(t *testing.T, store StringStore) {
	store.Set("mykey", "Hello world")
	assert.Equal(t, 11, store.Length("mykey"))
	assert.Equal(t, 0, store.Length("nonexisting"))
}
