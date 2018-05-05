package log

import (
	"context"
	"fmt"

	"google.golang.org/appengine/log"
)

type aeLogger struct {
}

// NewAELogger returns appengine logger
func NewAELogger() Logger {
	return &aeLogger{}
}

func (l *aeLogger) Debug(ctx context.Context, args ...interface{}) {
	l.Debugf(ctx, "%v", fmt.Sprint(args...))
}

func (l *aeLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	log.Debugf(ctx, format, args...)
}

func (l *aeLogger) Info(ctx context.Context, args ...interface{}) {
	l.Infof(ctx, "%v", fmt.Sprint(args...))
}

func (l *aeLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	log.Infof(ctx, format, args...)
}

func (l *aeLogger) Warning(ctx context.Context, args ...interface{}) {
	l.Warningf(ctx, "%v", fmt.Sprint(args...))
}

func (l *aeLogger) Warningf(ctx context.Context, format string, args ...interface{}) {
	log.Warningf(ctx, format, args...)
}

func (l *aeLogger) Error(ctx context.Context, args ...interface{}) {
	l.Errorf(ctx, "%v", fmt.Sprint(args...))
}

func (l *aeLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	log.Errorf(ctx, format, args...)
}

func (l *aeLogger) Critical(ctx context.Context, args ...interface{}) {
	l.Criticalf(ctx, "%v", fmt.Sprint(args...))
}

func (l *aeLogger) Criticalf(ctx context.Context, format string, args ...interface{}) {
	log.Criticalf(ctx, format, args...)
}
