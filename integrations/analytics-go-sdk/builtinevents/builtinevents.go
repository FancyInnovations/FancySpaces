package builtinevents

import (
	"log/slog"

	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/fancyinnovations/fancyspaces/analytics-sdk/client"
)

func ServiceStarted(service string) {
	evt := &client.Event{
		Name: "ServiceStarted",
		Properties: map[string]string{
			"service": service,
		},
	}

	if err := client.DefaultClient.SendEvent(evt); err != nil {
		slog.Warn("failed to send ServiceStarted event", sloki.WrapError(err))
	}
}

func ServiceStopped(service string) {
	evt := &client.Event{
		Name: "ServiceStopped",
		Properties: map[string]string{
			"service": service,
		},
	}

	if err := client.DefaultClient.SendEvent(evt); err != nil {
		slog.Warn("failed to send ServiceStopped event", sloki.WrapError(err))
	}
}
