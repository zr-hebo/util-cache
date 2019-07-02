package cache

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func Test_LRU(t *testing.T) {
	cache := NewLRUCache(40, 3)
	for i := 0; i < 1000; i++ {
		cache.Set(i, fmt.Sprintf("test data %d", i))
		//println(cache.String())
		time.Sleep(time.Millisecond * 50)
		a := i - int(rand.Int31n(500))
		if a < 0 {
			a = 0
		}

		val := cache.Get(a)
		if val == nil {
			fmt.Printf("get %d --> nil\n", a)
		} else {
			fmt.Printf("get %d --> %v\n", a, val)
		}
	}
}
