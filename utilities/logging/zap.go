package logging

import (
	"context"
	"fmt"
	"geolize/utilities/service"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapLogger implements the Logger interface using zap
type zapLogger struct {
	logger *zap.Logger
}

// NewZapLogger creates a new ZapLogger instance
func newZapLogger() (Logger, error) {
	var cores = make([]zapcore.Core, 0)

	var config zapcore.EncoderConfig
	if service.IsProd() {
		config = zap.NewProductionEncoderConfig()
	} else {
		config = zap.NewDevelopmentEncoderConfig()
	}

	config.EncodeTime = zapcore.RFC3339TimeEncoder

	if consoleLogEnabled {
		consoleEncoder := zapcore.NewConsoleEncoder(config)

		stdout := zapcore.AddSync(os.Stdout)
		cores = append(cores,
			zapcore.NewCore(consoleEncoder, stdout, getAtomicLevel(consoleLogLevel)),
		)
	}

	if logFileEnabled {
		jsonEncoder := zapcore.NewJSONEncoder(config)

		dir := filepath.Dir(logFilePath)

		// Ensure log directory exists
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		// Open log file
		logFile, err := os.OpenFile(logFilePath,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0644,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}

		// Create file core with async writer
		fileWriter := zapcore.AddSync(logFile)
		cores = append(cores,
			zapcore.NewCore(jsonEncoder, fileWriter, getAtomicLevel(logFileLevel)),
		)
	}

	// Create core with multiple writers
	core := zapcore.NewTee(cores...)

	// Create logger with the custom core
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return &zapLogger{
		logger: logger,
	}, nil
}

func getAtomicLevel(level string) zap.AtomicLevel {
	return zap.NewAtomicLevelAt(func() zapcore.Level {
		switch level {
		case "info":
			return zapcore.InfoLevel
		case "debug":
			return zapcore.DebugLevel
		case "error":
			return zapcore.ErrorLevel
		case "warn":
			return zapcore.WarnLevel
		case "fatal":
			return zapcore.FatalLevel
		default:
			return zapcore.DebugLevel
		}
	}())
}

// implement with Logger
func (z *zapLogger) Debug(ctx context.Context, msg string, keyvals ...KeyVal) {
	fields := make([]zap.Field, 0, len(keyvals))
	for _, kv := range keyvals {
		fields = append(fields, zap.Any(kv.Key, kv.Val))
	}
	z.logger.Debug(msg, fields...)
}

func (z *zapLogger) Info(ctx context.Context, msg string, keyvals ...KeyVal) {
	fields := make([]zap.Field, 0, len(keyvals))
	for _, kv := range keyvals {
		fields = append(fields, zap.Any(kv.Key, kv.Val))
	}
	z.logger.Info(msg, fields...)
}

func (z *zapLogger) Warn(ctx context.Context, msg string, keyvals ...KeyVal) {
	fields := make([]zap.Field, 0, len(keyvals))
	for _, kv := range keyvals {
		fields = append(fields, zap.Any(kv.Key, kv.Val))
	}
	z.logger.Warn(msg, fields...)
}

func (z *zapLogger) Error(ctx context.Context, msg string, keyvals ...KeyVal) {
	fields := make([]zap.Field, 0, len(keyvals))
	for _, kv := range keyvals {
		fields = append(fields, zap.Any(kv.Key, kv.Val))
	}
	z.logger.Error(msg, fields...)
}

func (z *zapLogger) Fatal(ctx context.Context, msg string, keyvals ...KeyVal) {
	fields := make([]zap.Field, 0, len(keyvals))
	for _, kv := range keyvals {
		fields = append(fields, zap.Any(kv.Key, kv.Val))
	}
	z.logger.Fatal(msg, fields...)
}

func (z *zapLogger) WithFields(keyvals ...KeyVal) Logger {
	fields := make([]zap.Field, 0, len(keyvals))
	for _, kv := range keyvals {
		fields = append(fields, zap.Any(kv.Key, kv.Val))
	}
	return &zapLogger{
		logger: z.logger.With(fields...),
	}
}

// Sync flushes any buffered log entries
func (z *zapLogger) Sync() error {
	return z.logger.Sync()
}
