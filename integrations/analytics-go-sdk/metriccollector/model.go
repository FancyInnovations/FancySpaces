package metriccollector

import "github.com/fancyinnovations/fancyspaces/integrations/analytics-go-sdk/client"

type MetricProvider func() ([]client.RecordData, error)
