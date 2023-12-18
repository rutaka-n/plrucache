package plrucache

import (
    "sync"
    "fmt"
	"testing"
	"time"
)

func TestStoreAndGet(t *testing.T) {
	t.Run("Store value and get it", func(t *testing.T) {
		cache := New[string](1, 1*time.Minute)
		key := "k1"
		value := "val1"
		cache.Set(key, value)
		res, ok := cache.Get(key)
		if !ok {
			t.Errorf("expect: %v, got: %v", true, ok)
		}
		if res != value {
			t.Errorf("expected: %s, got: %v", value, res)
		}
	})

	t.Run("Get value that does not exist", func(t *testing.T) {
		cache := New[string](1, 1*time.Minute)
		key := "k1"
		res, ok := cache.Get(key)
		if ok {
			t.Errorf("expect: %v, got: %v", false, ok)
		}
		if res != nil {
			t.Errorf("expected: %v, got: %v", nil, res)
		}
	})

	t.Run("Overfil the cache size and try to get displaced value", func(t *testing.T) {
		cache := New[string](1, 1*time.Minute)
		data := [][2]string{
			{"k1", "val1"},
			{"k2", "val2"},
		}
		for _, pair := range data {
			cache.Set(pair[0], pair[1])
		}
		// first pair should be displaced
		res, ok := cache.Get(data[0][0])
		if ok {
			t.Errorf("expect: %v, got: %v", false, ok)
		}
		if res != nil {
			t.Errorf("expected: %v, got: %v", nil, res)
		}
		// second pair should be in cache
		res, ok = cache.Get(data[1][0])
		if !ok {
			t.Errorf("expect: %v, got: %v", true, ok)
		}
		if res != data[1][1] {
			t.Errorf("expected: %s, got: %v", data[1][1], res)
		}
	})

	t.Run("Remove item from cache", func(t *testing.T) {
		cache := New[string](3, 1*time.Minute)
		data := [][2]string{
			{"k1", "val1"},
			{"k2", "val2"},
			{"k3", "val3"},
		}
		for _, pair := range data {
			cache.Set(pair[0], pair[1])
		}
		if cache.Len() != len(data) {
			t.Errorf("expected %d, got %d", len(data), cache.Len())
		}
		cache.Delete(data[0][0])
		if cache.Len() != len(data)-1 {
			t.Errorf("expected %d, got %d", len(data)-1, cache.Len())
		}
	})

	t.Run("Statistics calculation", func(t *testing.T) {
		cache := New[string](10, 1*time.Minute)
		data := [][2]string{
			{"k1", "val1"},
			{"k2", "val2"},
		}
		for _, pair := range data {
			cache.Set(pair[0], pair[1])
		}
		stat := cache.Stat()
		if stat.Misses != 0 {
			t.Errorf("expected: %d, got: %d", 0, stat.Misses)
		}
		if stat.Hits != 0 {
			t.Errorf("expected: %d, got: %d", 0, stat.Hits)
		}

		for _, pair := range data {
			for i := 0; i < 3; i++ {
				cache.Get(pair[0])
			}
		}

		stat = cache.Stat()
		if stat.Misses != 0 {
			t.Errorf("expected: %d, got: %d", 0, stat.Misses)
		}
		if stat.Hits != uint64(3*len(data)) {
			t.Errorf("expected: %d, got: %d", 3*len(data), stat.Hits)
		}

		for i := 0; i < 3; i++ {
			cache.Get("not exists")
		}

		stat = cache.Stat()
		if stat.Misses != 3 {
			t.Errorf("expected: %d, got: %d", 3, stat.Misses)
		}
		if stat.Hits != uint64(3*len(data)) {
			t.Errorf("expected: %d, got: %d", 3*len(data), stat.Hits)
		}

		cache.Reset()
        if cache.Len() != 0 {
            t.Errorf("expected: %v, got: %v", 0, cache.Len())
        }

		stat = cache.Stat()
		if stat.Misses != 0 {
			t.Errorf("expected: %d, got: %d", 0, stat.Misses)
		}
		if stat.Hits != 0 {
			t.Errorf("expected: %d, got: %d", 0, stat.Hits)
		}
	})
}

func benchmarkCache(size, i int, b *testing.B) {
		cache := New[string](size, 5*time.Second)
        data := make([][2]string, i)
        for j := 0; j < i; j++ {
            data[j][0] = fmt.Sprintf("k%d", j)
            data[j][1] = fmt.Sprintf("v%d", j)
        }

        b.ResetTimer()
        for n := 0; n < b.N; n++ {
            wg := sync.WaitGroup{}
            wg.Add(i*2)
            for j :=0; j < i; j++ {
                idx := j
                go func() {
                    idx := idx
                    cache.Set(data[idx][0], data[idx][1])
                    wg.Done()
                }()
                go func() {
                    idx := idx
                    _, _ = cache.Get(data[idx][0])
                    wg.Done()
                }()
            }
            wg.Wait()
        }
}

func BenchmarkCache100(b *testing.B)  { benchmarkCache(100, 100, b) }
func BenchmarkCache200(b *testing.B)  { benchmarkCache(100, 200, b) }
func BenchmarkCache300(b *testing.B)  { benchmarkCache(100, 300, b) }
func BenchmarkCache1000(b *testing.B) { benchmarkCache(100, 1000, b) }
func BenchmarkCache2000(b *testing.B) { benchmarkCache(100, 2000, b) }
func BenchmarkCache4000(b *testing.B) { benchmarkCache(100, 4000, b) }

func benchmarkCacheGet(size, i int, b *testing.B) {
		cache := New[string](size, 5*time.Second)
        data := make([][2]string, i)
        for j := 0; j < i; j++ {
            data[j][0] = fmt.Sprintf("k%d", j)
            data[j][1] = fmt.Sprintf("v%d", j)
            cache.Set(data[j][0], data[j][1])
        }

        b.ResetTimer()
        for n := 0; n < b.N; n++ {
            for j :=0; j < i; j++ {
                _, _ = cache.Get(data[j][0])
            }
        }
}

func BenchmarkCacheGet100(b *testing.B)  { benchmarkCacheGet(100, 100, b) }
func BenchmarkCacheGet200(b *testing.B)  { benchmarkCacheGet(100, 200, b) }
func BenchmarkCacheGet300(b *testing.B)  { benchmarkCacheGet(100, 300, b) }
func BenchmarkCacheGet1000(b *testing.B) { benchmarkCacheGet(100, 1000, b) }
func BenchmarkCacheGet2000(b *testing.B) { benchmarkCacheGet(100, 2000, b) }
func BenchmarkCacheGet4000(b *testing.B) { benchmarkCacheGet(100, 4000, b) }
