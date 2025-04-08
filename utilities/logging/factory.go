package logging

import (
	"fmt"
	"geolize/utilities/conf"
	"geolize/utilities/service"
)

// LoggerType represents the type of logger to create
type LoggerType string

const (
	// ZapLoggerType represents a Zap logger
	ZapLoggerType LoggerType = "zap"
)

var (
	consoleLogEnabled, _ = conf.GetBool("log_console", "enabled", true)
	consoleLogLevel, _   = conf.GetString("log_console", "level", "debug")

	logFileEnabled, _ = conf.GetBool("log_file", "enabled", true)
	logFileLevel, _   = conf.GetString("log_file", "level", "debug")
	logFilePath, _    = conf.GetString("log_file", "path", fmt.Sprintf("logs/%s.log", service.GetName()))
)

// NewLogger creates a new logger of the specified type
func NewLogger(loggerType LoggerType) (Logger, error) {
	switch loggerType {
	case ZapLoggerType:
		return newZapLogger()
	default:
		return nil, fmt.Errorf("unsupported logger type: %s", loggerType)
	}
}
