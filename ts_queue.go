package plrucache

import (
	"time"
)

// tsQ is a FIFO queue that tracks if any of items is expired.
// It expects that items in the queue added with timestamp that monotonically increasing.
type tsQ struct {
    q *staticQ
}

func newTSQ(size int) *tsQ {
    q := newQueue(size)
    return &tsQ{q: q}
}

// Push adds new item with timestamp.
func (q *tsQ) Push(val string, ts time.Time) int {
    return q.q.Push(val, ts)
}

// Pop drops and returns the item with oldest timestamp if exists.
func (q *tsQ) Pop() (qItem, bool) {
    return q.q.Pop()
}

// Delete item by index in the queue.
func (q *tsQ) Delete(idx int) bool {
    return q.q.Delete(idx)
}

// Len returns number of elements in the queue.
func (q *tsQ) Len() int {
    return q.q.Len()
}

// IsAnyExpired returns true if the next item from queue is expired.
func (q *tsQ) IsAnyExpired(now time.Time) bool {
    item, ok := q.q.Top()
    if !ok {
        return false
    }
	return now.After(item.ts)
}

