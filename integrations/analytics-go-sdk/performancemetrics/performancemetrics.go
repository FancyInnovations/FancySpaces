package performancemetrics

import (
	"runtime"
	"time"

	"github.com/fancyinnovations/fancyspaces/analytics-sdk/client"
)

type Service struct {
	appStartTime time.Time
	lastReset    time.Time
}

func NewService() *Service {
	s := &Service{
		appStartTime: time.Now(),
	}
	s.resetMetrics()

	return s
}

func (s *Service) resetMetrics() {
	s.lastReset = time.Now()
}

func (s *Service) collectAndReset() (*PerformanceMetrics, error) {
	cpuTime, err := getCpuTime()
	if err != nil {
		return nil, err
	}
	cpuUsage := cpuTime / time.Since(s.lastReset).Seconds()

	goroutineCount := runtime.NumGoroutine()

	allocedMemory := getAllocatedMemory()

	usedDisk, freeDisk, totalDisk, err := getDiskUsage("/")
	if err != nil {
		return nil, err
	}

	uptimeMs := time.Since(s.appStartTime).Milliseconds()

	metrics := &PerformanceMetrics{
		Uptime:          uptimeMs,
		CpuUsage:        cpuUsage,
		AllocatedMemory: allocedMemory,
		GoroutineCount:  goroutineCount,
		DiskUsage: DiskUsage{
			Used:  usedDisk,
			Free:  freeDisk,
			Total: totalDisk,
		},
	}

	s.resetMetrics()

	return metrics, nil
}

func (s *Service) MetricProvider() ([]client.RecordData, error) {
	metrics, err := s.collectAndReset()
	if err != nil {
		return nil, err
	}

	records := []client.RecordData{
		{
			Metric: "uptime",
			Label:  "",
			Value:  float64(metrics.Uptime),
		},
		{
			Metric: "cpu_usage",
			Label:  "",
			Value:  metrics.CpuUsage,
		},
		{
			Metric: "allocated_memory",
			Label:  "",
			Value:  float64(metrics.AllocatedMemory),
		},
		{
			Metric: "goroutine_count",
			Label:  "",
			Value:  float64(metrics.GoroutineCount),
		},
		{
			Metric: "disk_usage_total",
			Label:  "",
			Value:  float64(metrics.DiskUsage.Total),
		},
		{
			Metric: "disk_usage",
			Label:  "free",
			Value:  float64(metrics.DiskUsage.Free),
		},
		{
			Metric: "disk_usage",
			Label:  "used",
			Value:  float64(metrics.DiskUsage.Used),
		},
	}

	return records, nil
}
