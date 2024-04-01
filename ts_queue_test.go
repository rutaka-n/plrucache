package plrucache

import (
	"testing"
	"time"
)

func TestIsAnyExpired(t *testing.T) {
	q := newTSQ[string](1)
	val := "val"
	now := time.Now()
	q.Push(val, now.Add(10*time.Second))
	if got := q.IsAnyExpired(now); got {
		t.Errorf("expected: %v, got: %v", false, got)
	}
	if got := q.IsAnyExpired(now.Add(11 * time.Second)); !got {
		t.Errorf("expected: %v, got: %v", true, got)
	}
	if got := q.IsAnyExpired(now.Add(10 * time.Second)); !got {
		t.Errorf("expected: %v, got: %v", true, got)
	}
}
