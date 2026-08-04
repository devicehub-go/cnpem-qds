package queue

import "sync"

type Queue struct {
	mutex sync.Mutex
	items [][]byte

	Notify chan struct{}
}

// Instatiate a new queue
func New() *Queue {
	return &Queue{
		Notify: make(chan struct{}, 1),
	}
}

// Enqueue adds an item to the end
func (q *Queue) Enqueue(item []byte) {
	q.mutex.Lock()
	q.items = append(q.items, item)
	q.mutex.Unlock()

	select {
	case q.Notify <- struct{}{}:
	default:
		// Signal already pending, no need to add another
	}
}

// Dequeue removes and returns the first item
func (q *Queue) Dequeue() ([]byte, bool) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.items) == 0 {
		return nil, false
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

// Returns the number of pending items to be processed
func (q *Queue) Length() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	return len(q.items)
}
