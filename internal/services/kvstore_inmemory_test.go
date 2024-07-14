package services

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func BenchmarkInMemoryKeyValueStore_Set(b *testing.B) {
	store := NewInMemoryKeyValueStore()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = store.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}
}

func BenchmarkInMemoryKeyValueStore_Get(b *testing.B) {
	store := NewInMemoryKeyValueStore()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		_ = store.Set(key, value)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		actual, _ := store.Get(key)
		assert.Equal(b, value, actual)
	}
}

func BenchmarkInMemoryKeyValueStore_Delete(b *testing.B) {
	store := NewInMemoryKeyValueStore()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		_ = store.Set(key, value)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i)
		err := store.Delete(key)
		assert.NoError(b, err)
	}
}
