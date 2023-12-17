package plrucache

import (
	"testing"
	"time"
)

func TestQueuePushPop(t *testing.T) {
	t.Run("Pop from empty list", func(t *testing.T) {
		queue := newQueue[string](5)
		_, ok := queue.Pop()
		if ok {
			t.Errorf("expected: %v, got %v", false, ok)
		}
	})
	t.Run("Push and pop items", func(t *testing.T) {
		queue := newQueue[string](5)
		items := []string{"k1", "k2", "k3", "k4"}
		for _, k := range items {
			queue.Push(k, time.Now())
		}
		for i := len(items) - 1; i >= 0; i-- {
			_, ok := queue.Pop()
			if !ok {
				t.Fatalf("expected: %v, got %v", true, ok)
			}
		}
	})
}

func TestQueueDelete(t *testing.T) {
	t.Run("there is no elements", func(t *testing.T) {
		queue := newQueue[string](5)
		ok := queue.Delete(0)
		if ok {
			t.Fatalf("expected: %v, got: %v", false, ok)
		}
	})
	t.Run("delete head element", func(t *testing.T) {
		queue := newQueue[string](5)
		headIdx := queue.Push("head", time.Now()) // become head and tail
		_ = queue.Push("tail", time.Now())        // since it always pushes into the tail "tail" will be in the tail and "head" become head
		ok := queue.Delete(headIdx)
		if !ok {
			t.Fatalf("expected: %v, got: %v", true, ok)
		}
		// check that element actually was deleted
		if queue.Len() != 1 {
			t.Fatalf("expected: %d, got: %d", 1, queue.Len())
		}
		// check that last element is "tail"
		item, ok := queue.Pop()
		if !ok {
			t.Fatalf("expected: %v, got: %v", true, ok)
		}
		if item.val != "tail" {
			t.Fatalf("expected: %s, got: %s", "tail", item.val)
		}

	})
	t.Run("delete tail element", func(t *testing.T) {
		queue := newQueue[string](5)
		_ = queue.Push("head", time.Now())        // become head and tail
		tailIdx := queue.Push("tail", time.Now()) // since it always pushes into the tail "tail" will be in the tail and "head" become head
		ok := queue.Delete(tailIdx)
		if !ok {
			t.Fatalf("expected: %v, got: %v", true, ok)
		}
		// check that element actually was deleted
		if queue.Len() != 1 {
			t.Fatalf("expected: %d, got: %d", 1, queue.Len())
		}
		// check that last element is "head"
		item, ok := queue.Pop()
		if !ok {
			t.Fatalf("expected: %v, got: %v", true, ok)
		}
		if item.val != "head" {
			t.Fatalf("expected: %s, got: %s", "head", item.val)
		}

	})
	t.Run("delete element in the middle", func(t *testing.T) {
		queue := newQueue[string](5)
		_ = queue.Push("head", time.Now())
		midIdx := queue.Push("mid", time.Now())
		_ = queue.Push("tail", time.Now())
		ok := queue.Delete(midIdx)
		if !ok {
			t.Fatalf("expected: %v, got: %v", true, ok)
		}
		// check that element actually was deleted
		if queue.Len() != 2 {
			t.Fatalf("expected: %d, got: %d", 2, queue.Len())
		}
		// check that next element to pop is "head"
		item, ok := queue.Pop()
		if !ok {
			t.Fatalf("expected: %v, got: %v", true, ok)
		}
		if item.val != "head" {
			t.Fatalf("expected: %s, got: %s", "head", item.val)
		}
		// check that next element to pop is "tail"
		item, ok = queue.Pop()
		if !ok {
			t.Fatalf("expected: %v, got: %v", true, ok)
		}
		if item.val != "tail" {
			t.Fatalf("expected: %s, got: %s", "tail", item.val)
		}
	})
}
