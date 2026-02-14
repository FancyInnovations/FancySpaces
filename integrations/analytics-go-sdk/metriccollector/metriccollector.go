package metriccollector

import (
	"log/slog"
	"time"

	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/analytics-sdk/client"
)

type Service struct {
	c                *client.Client
	interval         time.Duration
	providers        []MetricProvider
	schedulerStarted bool
	abortScheduler   chan struct{}
}

type Configuration struct {
	Client    *client.Client
	Interval  time.Duration
	Providers []MetricProvider
}

func NewService(cfg *Configuration) *Service {
	if cfg.Interval <= 0 {
		cfg.Interval = 60
	}

	if cfg.Providers == nil {
		cfg.Providers = []MetricProvider{}
	}

	return &Service{
		c:                cfg.Client,
		interval:         cfg.Interval,
		providers:        cfg.Providers,
		schedulerStarted: false,
		abortScheduler:   make(chan struct{}),
	}
}

func (s *Service) AddProvider(provider MetricProvider) {
	s.providers = append(s.providers, provider)
}

func (s *Service) Send() error {
	var records []client.RecordData
	for _, provider := range s.providers {
		providerRecords, err := provider()
		if err != nil {
			continue
		}
		records = append(records, providerRecords...)
	}

	if err := s.c.SendRecord(records); err != nil {
		return err
	}

	return nil
}

func (s *Service) StartScheduler() {
	if s.schedulerStarted {
		return
	}

	s.schedulerStarted = true

	go func() {
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := s.Send(); err != nil {
					slog.Warn("failed to send metrics", sloki.WrapError(err))
				}
			case <-s.abortScheduler:
				return
			}
		}
	}()
}

func (s *Service) StopScheduler() {
	s.abortScheduler <- struct{}{}
	s.schedulerStarted = false
}
