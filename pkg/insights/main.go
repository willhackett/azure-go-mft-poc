package insights

import (
	"fmt"
	"time"

	"github.com/microsoft/ApplicationInsights-Go/appinsights"
	"github.com/willhackett/azure-mft/pkg/config"
)

var (
	client appinsights.TelemetryClient
)

func Init() {
	instrumentationKey := config.GetConfig().Azure.InstrumentationKey
	if instrumentationKey == "" {
		fmt.Println("No telemetry key provided so only console output will be available.")
		return
	}

	telemetryConfig := appinsights.NewTelemetryConfiguration(instrumentationKey)
	telemetryConfig.MaxBatchSize = 1024
	telemetryConfig.MaxBatchInterval = 2 * time.Second
	client = appinsights.NewTelemetryClientFromConfig(telemetryConfig)

	client.TrackEvent("Application Startup")
}
