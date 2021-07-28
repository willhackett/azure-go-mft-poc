package insights

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/microsoft/ApplicationInsights-Go/appinsights"
	"github.com/microsoft/ApplicationInsights-Go/appinsights/contracts"
	"github.com/sirupsen/logrus"
	"github.com/willhackett/azure-mft/pkg/config"
)

var (
	client appinsights.TelemetryClient

	logLevelMap = map[logrus.Level]contracts.SeverityLevel{
		logrus.FatalLevel: appinsights.Critical,
		logrus.ErrorLevel: appinsights.Error,
		logrus.WarnLevel:  appinsights.Warning,
		logrus.InfoLevel:  appinsights.Information,
	}
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
}

// InsightsHook is a logrus hook
type InsightsHook struct{}

// Levels returns logging level to fire this hook.
func (hook *InsightsHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
	}
}

// Fire is called by logurs
func (hook *InsightsHook) Fire(entry *logrus.Entry) error {
	if client == nil {
		return nil
	}

	go hook.fire(entry)
	return nil
}

func (hook *InsightsHook) fire(entry *logrus.Entry) {
	trace, _ := hook.buildTrace(entry)
	client.Track(trace)
}

func (hook *InsightsHook) buildTrace(entry *logrus.Entry) (*appinsights.TraceTelemetry, error) {
	if client == nil {
		return nil, nil
	}

	if _, ok := entry.Data["message"]; !ok {
		entry.Data["message"] = entry.Message
	}

	level := logLevelMap[entry.Level]
	trace := appinsights.NewTraceTelemetry(entry.Message, level)

	if trace == nil {
		return nil, errors.New("failed to create trace")
	}

	trace.Properties["log_level"] = entry.Level.String()

	for k, v := range entry.Data {
		v = formatData(v)
		trace.Properties[k] = fmt.Sprintf("%v", v)
	}
	return trace, nil
}

// formatData returns value as a suitable format.
func formatData(value interface{}) (formatted interface{}) {
	switch value := value.(type) {
	case json.Marshaler:
		return value
	case error:
		return value.Error()
	case fmt.Stringer:
		return value.String()
	default:
		return value
	}
}
