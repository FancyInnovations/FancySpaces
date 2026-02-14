package performancemetrics

import "runtime"

func getAllocatedMemory() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc
}
