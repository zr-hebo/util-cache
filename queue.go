package cache

import (
	"sync"
)

// Node
type node struct {
	value interface{}
	next  *node
	pre  *node
}

type LinkedQueue struct {
	head  *node
	tail  *node
	size  int
	mutex sync.Mutex
}

func NewLinkedQueue() (lq *LinkedQueue) {
	lq = &LinkedQueue{
		head: &node{},
	}
	return
}

func (lq *LinkedQueue) Enqueue(val interface{}) {
	lq.mutex.Lock()
	defer lq.mutex.Unlock()

	newNode := &node{
		value: val,
		next:  lq.head,
	}
	newNode.next = lq.head.next
	lq.head.next = newNode

	lq.size = lq.size + 1
}

func (lq *LinkedQueue) Dequeue() (val interface{}) {
	lq.mutex.Lock()
	defer lq.mutex.Unlock()

	if lq.head == nil {
		return nil
	}

	val = lq.head.value
	lq.head = lq.head.next

	lq.size = lq.size - 1
	return
}
