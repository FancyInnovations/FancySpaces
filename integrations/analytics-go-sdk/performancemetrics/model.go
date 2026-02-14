package performancemetrics

type PerformanceMetrics struct {
	Uptime          int64
	CpuUsage        float64
	AllocatedMemory uint64
	GoroutineCount  int
	DiskUsage       DiskUsage
}

type DiskUsage struct {
	Used  uint64
	Free  uint64
	Total uint64
}
