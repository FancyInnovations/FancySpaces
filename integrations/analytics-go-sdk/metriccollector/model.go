package metriccollector

import "github.com/fancyinnovations/fancyspaces/analytics-sdk/client"

type MetricProvider func() ([]client.RecordData, error)
