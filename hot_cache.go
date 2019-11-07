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

func (r *Record) String() string {
	return fmt.Sprintf(
		"key:%v, last visit time:%s", r.listIdx.Value,
		r.lastVisitTime.Format("2006-01-02 15:04:05"))
}

// Record record
type Record struct {
	listIdx       *list.Element
	val           interface{}
	lastVisitTime time.Time
}

// LRUCache lru cache
type LRUCache struct {
	// maxNum is the maximum number of cache entries before
	maxNum           int
	lock             *sync.RWMutex
	keyOrder         *list.List
	contents         map[Key]*Record
	timeoutInSeconds int
}

// NewLRUCache create LRUCache instance
func NewLRUCache(num, timeoutInSeconds int) (rc *LRUCache) {
	if num < 1 {
		num = 0
	}

	return &LRUCache{
		maxNum:           num,
		lock:             &sync.RWMutex{},
		keyOrder:         list.New(),
		contents:         make(map[Key]*Record),
		timeoutInSeconds: timeoutInSeconds,
	}
}

func (lc LRUCache) String() string {
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
		lc.keyOrder.MoveToFront(record.listIdx)

	} else {
		// add new element
		if lc.keyOrder.Len() >= lc.maxNum {
			leastUsedElement := lc.keyOrder.Back()
			delete(lc.contents, key)
			lc.keyOrder.Remove(leastUsedElement)
		}

		record = new(Record)
		record.val = val
		record.listIdx = lc.keyOrder.PushFront(key)
		lc.contents[key] = record
	}

	record.lastVisitTime = time.Now()
}

// Get get value by key
func (lc *LRUCache) Get(key interface{}) (val interface{}) {
	lc.lock.RLock()
	defer lc.lock.RUnlock()

	record, ok := lc.contents[key]
	if ok {
		// check if is expired record
		now := time.Now()
		behindSeconds := time.Second * time.Duration(lc.timeoutInSeconds)
		if now.After(record.lastVisitTime.Add(behindSeconds)) {
			lc.keyOrder.Remove(record.listIdx)
			delete(lc.contents, key)
			return nil
		}

		// reorder key position
		lc.keyOrder.PushFront(key)
		val = record.val
		record.lastVisitTime = now
		return val
	}
	return nil
}

// Remove remove value by key
func (lc *LRUCache) Remove(key interface{}) (val interface{}) {
	lc.lock.Lock()
	defer lc.lock.Unlock()

	record, ok := lc.contents[key]
	if ok {
		lc.keyOrder.Remove(record.listIdx)
		delete(lc.contents, key)
	}
	return
}
