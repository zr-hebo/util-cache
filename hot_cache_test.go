package cache

import (
	"fmt"
	"testing"
	"time"
)

func Test_LRU(t *testing.T) {
	cache := NewLRUCache(40, 3)
	for i := 0; i < 1000; i++ {
		cache.Set(i, fmt.Sprintf("test data %d", i))
		println(cache.String())
		time.Sleep(time.Millisecond * 10)
	}
}
