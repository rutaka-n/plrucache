package plrucache

import (
	"time"
)

// tsQ is a FIFO queue that tracks if any of items is expired.
// It expects that items in the queue added with timestamp that monotonically increasing.
type tsQ[T any] struct {
    q *staticQ[T]
}

func newTSQ[T any](size int) *tsQ[T] {
    q := newQueue[T](size)
    return &tsQ[T]{q: q}
}

// Push adds new item with timestamp.
func (q *tsQ[T]) Push(val T, ts time.Time) int {
    return q.q.Push(val, ts)
}

// Pop drops and returns the item with oldest timestamp if exists.
func (q *tsQ[T]) Pop() (qItem[T], bool) {
    return q.q.Pop()
}

// Delete item by index in the queue.
func (q *tsQ[T]) Delete(idx int) bool {
    return q.q.Delete(idx)
}

// Len returns number of elements in the queue.
func (q *tsQ[T]) Len() int {
    return q.q.Len()
}

// IsAnyExpired returns true if the next item from queue is expired.
func (q *tsQ[T]) IsAnyExpired(now time.Time) bool {
    item, ok := q.q.Top()
    if !ok {
        return false
    }
	return now.After(item.ts)
}

