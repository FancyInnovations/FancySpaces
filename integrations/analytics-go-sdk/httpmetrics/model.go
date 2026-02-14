package httpmetrics

import "net/http"

type StatusRecorder struct {
	http.ResponseWriter
	Status int
	Size   int64
}

func (s *StatusRecorder) WriteHeader(code int) {
	s.Status = code
	s.ResponseWriter.WriteHeader(code)
}

func (s *StatusRecorder) Write(b []byte) (int, error) {
	n, err := s.ResponseWriter.Write(b)
	s.Size += int64(n)
	return n, err
}

type RequestMetrics struct {
	RequestsPerSecond float64
	StatusCodes       map[string]int64
	Durations         Stats
	RequestSizes      Stats
	ResponseSizes     Stats
}

type Stats struct {
	Min int64
	Max int64
	Avg float64
}
