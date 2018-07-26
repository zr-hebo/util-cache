package cache

import (
	"bytes"
	"container/list"
	"fmt"
	"sync"
	"time"
)

// Key key
type Key interface{}

// Record record
type Record struct {
	key           *list.Element
	val           interface{}
	lastVisitTime time.Time
}

// LRUCache lru cache
type LRUCache struct {
	// maxNum is the maximum number of cache entries before
	maxNum      int
	lock        *sync.Mutex
	keyOrder    *list.List
	contents    map[Key]*Record
	keepSeconds int
}

// NewLRUCache create LRUCache instance
func NewLRUCache(num, keepSeconds int) (rc *LRUCache) {
	if num < 1 {
		num = 0
	}

	return &LRUCache{
		maxNum:      num,
		lock:        &sync.Mutex{},
		keyOrder:    list.New(),
		contents:    make(map[Key]*Record),
		keepSeconds: keepSeconds,
	}
}

func (r *Record) String() string {
	return fmt.Sprintf(
		"key:%s, last visit time:%s", r.key.Value,
		r.lastVisitTime.Format("%Y-%m-%d %H:%M:%S"))
}

func (lc *LRUCache) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "there are %d record in LRU cache. ", lc.keyOrder.Len())
	for _, record := range lc.contents {
		fmt.Fprint(&buf, fmt.Sprintf("%s; ", record))
	}

	return buf.String()
}

// Set add kv item to lru cache
func (lc *LRUCache) Set(key, val interface{}) {
	lc.lock.Lock()
	defer lc.lock.Unlock()

	record, ok := lc.contents[key]
	if ok {
		// update value in element
		record.val = val
		lc.keyOrder.MoveToBack(record.key)

	} else {
		// add new element
		if lc.keyOrder.Len() >= lc.maxNum {
			leastUsedElement := lc.keyOrder.Front()
			lc.keyOrder.Remove(leastUsedElement)
			delete(lc.contents, leastUsedElement.Value)
		}

		record = new(Record)
		record.val = val
		record.key = lc.keyOrder.PushBack(key)
		lc.contents[key] = record
	}

	record.lastVisitTime = time.Now()
}

// Get get value by key
func (lc *LRUCache) Get(key interface{}) (val interface{}) {
	isExpired := false
	defer func() {
		if isExpired {
			lc.Remove(key)
		}
	}()

	lc.lock.Lock()
	defer lc.lock.Unlock()

	record, ok := lc.contents[key]
	if ok {
		// check if is expired record
		now := time.Now()
		behindSeconds := time.Second * time.Duration(lc.keepSeconds)
		if now.After(record.lastVisitTime.Add(behindSeconds)) {
			isExpired = true
			return
		}

		// reorder key position
		lc.keyOrder.PushBack(record.key)
		val = record.val
		record.lastVisitTime = time.Now()
	}
	return
}

// Remove remove value by key
func (lc *LRUCache) Remove(key interface{}) (val interface{}) {
	lc.lock.Lock()
	defer lc.lock.Unlock()

	record, ok := lc.contents[key]
	if ok {
		lc.keyOrder.Remove(record.key)
		delete(lc.contents, key)
	}
	return
}
