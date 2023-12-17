package plrucache

import (
	"time"
)

// qItem contains value of queue item and timestamp.
type qItem[T any] struct {
	// according to the docs time.Time conatins monotonic time, so
	// it easy use time.Time to track timestampts instead of unixtime.
	ts   time.Time
	val  T
	prev int
	next int
}

type staticQ[T any] struct {
	maxSize   int
	head      int
	tail      int
	slots     []qItem[T]
	freeSlots map[int]interface{}
}

// newQueue returns empty statis queue with defined size.
// Under the hood it used fixed size slice that implement double-linked list queue,
// that allows to pop least recentrly used for O(1), push new key for O(1) and hit the existing key
// for O(1) in case the caller has pointer to the assosiated item.
func newQueue[T any](size int) *staticQ[T] {
	slots := make([]qItem[T], size)
	freeSlots := make(map[int]interface{}, size)
	for i := range slots {
		freeSlots[i] = nil
	}
	return &staticQ[T]{
		maxSize:   size,
		head:      -1,
		tail:      -1,
		slots:     slots,
		freeSlots: freeSlots,
	}
}

// Push inserts new item with specified timestamp in the tail of the list.
// Important: for proper work of queue timestamps should monotonically increasing for every new item.
func (q *staticQ[T]) Push(val T, ts time.Time) int {
	idx, ok := q.getFreeSlot()
	if !ok {
		// if no free slots left, pop least recently used
		// and accure freed up one
		q.Pop()
		idx, _ = q.getFreeSlot()
	}
	q.slots[idx] = qItem[T]{
		val:  val,
		ts:   time.Now(),
		prev: -1,
		next: -1,
	}
	// if tail does not exist, head is also does not exist,
	// so just put the first element into the list
	if q.tail == -1 {
		q.tail = idx
		q.head = idx
		return idx
	}
	oldTailIdx := q.tail
	q.tail = idx
	q.slots[idx].prev = oldTailIdx
	q.slots[oldTailIdx].next = idx
	return idx
}

// Pop removes and returns item from head of the list if there is one.
func (q *staticQ[T]) Pop() (qItem[T], bool) {
	if q.head == -1 {
		return qItem[T]{}, false
	}
	item := q.slots[q.head]
	q.freeSlots[q.head] = nil // mark slots as free
	if q.head == q.tail {
		q.head = -1
		q.tail = -1
	} else {
		q.slots[item.next].prev = -1
		q.head = item.next
	}
	return item, true
}

// Delete key if it is exist.
func (q *staticQ[T]) Delete(idx int) bool {
	if _, ok := q.freeSlots[idx]; ok {
		return false
	}
	prev := q.slots[idx].prev
	next := q.slots[idx].next
	if prev >= 0 {
		q.slots[prev].next = next
	} else {
		// it was the head item
		q.head = next
	}
	if next >= 0 {
		q.slots[next].prev = prev
	} else {
		// it was the tail item
		q.tail = prev
	}
	q.freeSlots[idx] = nil
	return true
}

func (q *staticQ[T]) Len() int {
	return q.maxSize - len(q.freeSlots)
}

// Top returns item from head of the list if there is one.
func (q *staticQ[T]) Top() (qItem[T], bool) {
	if q.head == -1 {
		return qItem[T]{}, false
	}
	item := q.slots[q.head]
	return item, true
}

// getFreeSlot returns random free slot in the slots slice, or (0, false) if no slots are available
func (q *staticQ[T]) getFreeSlot() (int, bool) {
	if len(q.freeSlots) == 0 {
		return 0, false
	}
	var idx int
	// get any item from set
	for idx = range q.freeSlots {
		break
	}
	delete(q.freeSlots, idx)
	return idx, true
}
