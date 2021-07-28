package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/willhackett/azure-mft/pkg/config"
	"github.com/willhackett/azure-mft/pkg/insights"
)

var (
	app = "startup"

	hostname = "unknown"
)

func init() {
	hostname, _ = os.Hostname()
}

func Init() {
	log.SetFormatter(&log.TextFormatter{})

	log.SetOutput(os.Stdout)

	var logLevel log.Level
	switch config.GetConfig().Agent.LogLevel {
	case "debug":
		logLevel = log.DebugLevel
	default:
		logLevel = log.InfoLevel
	}

	insightsHook := &insights.InsightsHook{}

	log.SetLevel(logLevel)
	log.AddHook(insightsHook)
}

func SetApp(a string) {
	app = a
}

func Get() *log.Entry {
	return log.WithFields(log.Fields{
		"app":      app,
		"agent":    config.GetConfig().Agent.Name,
		"hostname": hostname,
	})
}
