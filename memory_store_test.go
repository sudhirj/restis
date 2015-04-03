package restis

import "testing"

func TestMemoryStore(t *testing.T) {
	storeGenerator := func() Store {
		return NewMemoryStore()
	}
	RunAllTestsOnStore(t, storeGenerator)
	RunAllRedisDocChecksOnStore(t, storeGenerator)
}
