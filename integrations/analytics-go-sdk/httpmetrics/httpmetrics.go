package httpmetrics

import (
	"net/http"
	"sync"
	"time"

	"github.com/fancyinnovations/fancyspaces/integrations/analytics-go-sdk/client"
)

type Service struct {
	lastReset     time.Time
	statusCodes   map[string]int64
	requestSizes  []int64
	responseSizes []int64
	durations     []int64
	mu            sync.Mutex
}

func NewService() *Service {
	s := &Service{}
	s.resetMetrics()

	return s
}

func (s *Service) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		sr := &StatusRecorder{
			ResponseWriter: w,
			Status:         http.StatusOK,
		}

		next.ServeHTTP(sr, r)

		duration := time.Since(startTime).Milliseconds()

		s.mu.Lock()
		defer s.mu.Unlock()
		s.durations = append(s.durations, duration)
		s.requestSizes = append(s.requestSizes, r.ContentLength)
		s.responseSizes = append(s.responseSizes, sr.Size)

		statusCodeCategory := sr.Status / 100
		switch statusCodeCategory {
		case 2:
			s.statusCodes["2xx"]++
		case 3:
			s.statusCodes["3xx"]++
		case 4:
			s.statusCodes["4xx"]++
		case 5:
			s.statusCodes["5xx"]++
		}
	})
}

func (s *Service) collectAndReset() RequestMetrics {
	s.mu.Lock()
	defer s.mu.Unlock()

	// duration stats
	var totalDuration int64
	var minDuration int64 = -1
	var maxDuration int64 = -1

	for _, d := range s.durations {
		totalDuration += d
		if minDuration == -1 || d < minDuration {
			minDuration = d
		}
		if d > maxDuration {
			maxDuration = d
		}
	}

	var avgDuration float64
	if len(s.durations) > 0 {
		avgDuration = float64(totalDuration) / float64(len(s.durations))
	}

	// request size stats
	var totalRequestSize int64
	var minRequestSize int64 = -1
	var maxRequestSize int64 = -1

	for _, rs := range s.requestSizes {
		totalRequestSize += rs
		if minRequestSize == -1 || rs < minRequestSize {
			minRequestSize = rs
		}
		if rs > maxRequestSize {
			maxRequestSize = rs
		}
	}

	var avgRequestSize float64
	if len(s.requestSizes) > 0 {
		avgRequestSize = float64(totalRequestSize) / float64(len(s.requestSizes))
	}

	// response size stats
	var totalResponseSize int64
	var minResponseSize int64 = -1
	var maxResponseSize int64 = -1

	for _, rs := range s.responseSizes {
		totalResponseSize += rs
		if minResponseSize == -1 || rs < minResponseSize {
			minResponseSize = rs
		}
		if rs > maxResponseSize {
			maxResponseSize = rs
		}
	}

	var avgResponseSize float64
	if len(s.responseSizes) > 0 {
		avgResponseSize = float64(totalResponseSize) / float64(len(s.responseSizes))
	}

	// requests per second
	rps := float64(len(s.durations)) / time.Since(s.lastReset).Seconds()

	metrics := RequestMetrics{
		RequestsPerSecond: rps,
		StatusCodes:       s.statusCodes,
	}
	metrics.Durations.Min = minDuration
	metrics.Durations.Max = maxDuration
	metrics.Durations.Avg = avgDuration
	metrics.RequestSizes.Min = minRequestSize
	metrics.RequestSizes.Max = maxRequestSize
	metrics.RequestSizes.Avg = avgRequestSize
	metrics.ResponseSizes.Min = minResponseSize
	metrics.ResponseSizes.Max = maxResponseSize
	metrics.ResponseSizes.Avg = avgResponseSize

	s.resetMetrics()

	return metrics
}

func (s *Service) resetMetrics() {
	s.statusCodes = map[string]int64{
		"2xx": 0,
		"3xx": 0,
		"4xx": 0,
		"5xx": 0,
	}
	s.durations = []int64{}
	s.requestSizes = []int64{}
	s.responseSizes = []int64{}
	s.lastReset = time.Now()
}

func (s *Service) MetricProvider() ([]client.RecordData, error) {
	metrics := s.collectAndReset()

	records := []client.RecordData{
		{
			Metric: "requests_per_second",
			Label:  "",
			Value:  metrics.RequestsPerSecond,
		},
	}

	records = append(records,
		client.RecordData{
			Metric: "min_request_duration",
			Label:  "",
			Value:  float64(metrics.Durations.Min),
		},
		client.RecordData{
			Metric: "max_request_duration",
			Label:  "",
			Value:  float64(metrics.Durations.Max),
		},
		client.RecordData{
			Metric: "avg_request_duration",
			Label:  "",
			Value:  metrics.Durations.Avg,
		},
	)

	records = append(records,
		client.RecordData{
			Metric: "min_request_size",
			Label:  "",
			Value:  float64(metrics.RequestSizes.Min),
		},
		client.RecordData{
			Metric: "max_request_size",
			Label:  "",
			Value:  float64(metrics.RequestSizes.Max),
		},
		client.RecordData{
			Metric: "avg_request_size",
			Label:  "",
			Value:  metrics.RequestSizes.Avg,
		},
	)

	records = append(records,
		client.RecordData{
			Metric: "min_response_size",
			Label:  "",
			Value:  float64(metrics.ResponseSizes.Min),
		},
		client.RecordData{
			Metric: "max_response_size",
			Label:  "",
			Value:  float64(metrics.ResponseSizes.Max),
		},
		client.RecordData{
			Metric: "avg_response_size",
			Label:  "",
			Value:  metrics.ResponseSizes.Avg,
		},
	)

	for status, count := range metrics.StatusCodes {
		records = append(records, client.RecordData{
			Metric: "status_code_count",
			Label:  status,
			Value:  float64(count),
		})
	}

	return records, nil
}
