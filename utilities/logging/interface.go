package logging

import (
	"context"
)

type Logger interface {
	Debug(ctx context.Context, msg string, keyvals ...KeyVal)
	Info(ctx context.Context, msg string, keyvals ...KeyVal)
	Warn(ctx context.Context, msg string, keyvals ...KeyVal)
	Error(ctx context.Context, msg string, keyvals ...KeyVal)
	Fatal(ctx context.Context, msg string, keyvals ...KeyVal)

	WithFields(keyvals ...KeyVal) Logger
}
