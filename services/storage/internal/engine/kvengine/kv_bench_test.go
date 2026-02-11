package kvengine

import (
	"fmt"
	"testing"
	"time"

	"github.com/fancyinnovations/fancyspaces/storage/pkg/codex"
)

func newTestEngine() *Engine {
	return NewEngine(Configuration{
		DisableTTL: true,
	})
}

func testValue() *codex.Value {
	val, _ := codex.NewValue("Hello, World!")
	return val
}

func BenchmarkEngine_Set(b *testing.B) {
	e := newTestEngine()
	val := testValue()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Set(fmt.Sprintf("key-%d", i), val)
	}
}

func BenchmarkEngine_Get(b *testing.B) {
	e := newTestEngine()
	val := testValue()

	for i := 0; i < b.N; i++ {
		e.Set(fmt.Sprintf("key-%d", i), val)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Get(fmt.Sprintf("key-%d", i))
	}
}

func BenchmarkEngine_SetWithTTL(b *testing.B) {
	e := NewEngine(Configuration{DisableTTL: false})
	val := testValue()
	expires := time.Now().Add(1 * time.Minute).UnixNano()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.SetWithTTL(fmt.Sprintf("key-%d", i), val, expires)
	}
}

func BenchmarkEngine_GetExpired(b *testing.B) {
	e := NewEngine(Configuration{DisableTTL: false})
	val := testValue()
	expires := time.Now().Add(-1 * time.Second).UnixNano()

	for i := 0; i < b.N; i++ {
		e.SetWithTTL(fmt.Sprintf("key-%d", i), val, expires)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Get(fmt.Sprintf("key-%d", i))
	}
}

func BenchmarkEngine_SetMultiple(b *testing.B) {
	e := newTestEngine()

	entries := make(map[string]codex.Value, 100)
	for i := 0; i < 100; i++ {
		val, _ := codex.NewValue("value")
		entries[fmt.Sprintf("key-%d", i)] = *val
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.SetMultiple(entries)
	}
}

func BenchmarkEngine_GetMultiple(b *testing.B) {
	e := newTestEngine()

	keys := make([]string, 100)
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key-%d", i)
		keys[i] = key
		e.Set(key, testValue())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.GetMultiple(keys)
	}
}

func BenchmarkEngine_Delete(b *testing.B) {
	e := newTestEngine()

	for i := 0; i < b.N; i++ {
		e.Set(fmt.Sprintf("key-%d", i), testValue())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Delete(fmt.Sprintf("key-%d", i))
	}
}

func BenchmarkEngine_Exists(b *testing.B) {
	e := newTestEngine()
	val := testValue()

	for i := 0; i < b.N; i++ {
		e.Set(fmt.Sprintf("key-%d", i), val)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Exists(fmt.Sprintf("key-%d", i))
	}
}

func BenchmarkEngine_Size(b *testing.B) {
	e := newTestEngine()

	for i := 0; i < 100_000; i++ {
		e.Set(fmt.Sprintf("key-%d", i), testValue())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Count()
	}
}

func BenchmarkEngine_ParallelGet(b *testing.B) {
	e := newTestEngine()
	val := testValue()

	for i := 0; i < 100_000; i++ {
		e.Set(fmt.Sprintf("key-%d", i), val)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			_ = e.Get(fmt.Sprintf("key-%d", i%100_000))
			i++
		}
	})
}

func BenchmarkEngine_ParallelSet(b *testing.B) {
	e := newTestEngine()
	val := testValue()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			e.Set(fmt.Sprintf("key-%d", i), val)
			i++
		}
	})
}
