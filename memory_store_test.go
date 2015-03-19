package restis

import "testing"

func TestMemoryStore(t *testing.T) {
	memoryStore := NewMemoryStore()
	RunAllTestsOnStore(t, memoryStore)
}
