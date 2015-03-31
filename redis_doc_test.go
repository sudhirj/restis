package restis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func RunAllRedisDocChecksOnStore(t *testing.T, store Store) {
	APPEND(t, store)
	DECR(t, store)
	DECRBY(t, store)
	GET(t, store)
	GETRANGE(t, store)
	SET(t, store)

	INCR(t, store)
	INCRBY(t, store)
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
	// assert.Equal(t, "This", store.GetRange("mykey", 0, 3))
}

// redis> GETRANGE mykey 0 3
// "This"
// redis> GETRANGE mykey -3 -1
// "ing"
// redis> GETRANGE mykey 0 -1
// "This is a string"
// redis> GETRANGE mykey 10 100
// "string"

func SET(t *testing.T, store StringStore) {
	store.Set("mykey", "Hello")
	assert.Equal(t, "Hello", store.Get("mykey"))
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
