package restis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func CheckGET(t *testing.T, store StringStore) {
	assert.Equal(t, "", store.Get("nonexisting"))
	store.Set("mykey", "Hello")
	assert.Equal(t, "Hello", store.Get("mykey"))
}

func CheckSET(t *testing.T, store StringStore) {
	store.Set("mykey", "Hello")
	assert.Equal(t, "Hello", store.Get("mykey"))
}

func CheckDECR(t *testing.T, store StringStore) {
	store.Set("mykey", "10")
	assert.Equal(t, 9, store.Decrement("mykey"))
	store.Set("mykey", "234293482390480948029348230948")
	assert.Equal(t, -1, store.Decrement("mykey")) // DEVIATION: Redis throws an error instead.
}

func CheckDECRBY(t *testing.T, store StringStore) {
	store.Set("mykey", "10")
	assert.Equal(t, 7, store.DecrementBy("mykey", 3))
}

func CheckINCR(t *testing.T, store StringStore) {
	store.Set("mykey", "10")
	assert.Equal(t, 11, store.Increment("mykey"))
	assert.Equal(t, "11", store.Get("mykey"))
}

func CheckINCRBY(t *testing.T, store StringStore) {
	store.Set("mykey", "10")
	assert.Equal(t, 15, store.IncrementBy("mykey", 5))
}

func CheckGETRANGE(t *testing.T, store StringStore) {
	store.Set("mykey", "This is a string")
	assert.Equal(t, "This", store.GetRange("mykey", 0, 3))
	assert.Equal(t, "ing", store.GetRange("mykey", -3, -1))
	assert.Equal(t, "This is a string", store.GetRange("mykey", 0, -1))
	assert.Equal(t, "string", store.GetRange("mykey", 10, 100))
}

func RunAllRedisDocChecksOnStore(t *testing.T, store Store) {
	CheckGET(t, store)
	CheckSET(t, store)
	CheckDECR(t, store)
	CheckDECRBY(t, store)
	CheckINCR(t, store)
	CheckINCRBY(t, store)
	CheckGETRANGE(t, store)
}
