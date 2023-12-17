package plrucache

import (
	"time"
)

type item[T any] struct {
	expiration time.Time
	val        T
	lruIdx     int
	tsqIdx     int
}

func (i item[T]) IsExpired(now time.Time) bool {
	return now.After(i.expiration)
}
