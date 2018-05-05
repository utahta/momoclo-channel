package log

import "context"

// Logger interface
type Logger interface {
	Debug(ctx context.Context, args ...interface{})
	Debugf(ctx context.Context, format string, args ...interface{})

	Info(ctx context.Context, args ...interface{})
	Infof(ctx context.Context, format string, args ...interface{})

	Warning(ctx context.Context, args ...interface{})
	Warningf(ctx context.Context, format string, args ...interface{})

	Error(ctx context.Context, args ...interface{})
	Errorf(ctx context.Context, format string, args ...interface{})

	Critical(ctx context.Context, args ...interface{})
	Criticalf(ctx context.Context, format string, args ...interface{})
}
