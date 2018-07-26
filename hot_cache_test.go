package cache

import (
	"fmt"
	"testing"
	"time"
)

func Test_LRU(t *testing.T) {
	cache := NewLRUCache(100, 3)
	for i := 0; i < 10000; i++ {
		cache.Set(i, fmt.Sprintf("test data %d", i))
		print(cache)
		time.Sleep(1)
	}
}
