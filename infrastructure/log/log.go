package log

import (
	"context"
	"fmt"

	"github.com/utahta/momoclo-channel/domain/core"
	"google.golang.org/appengine/log"
)

type aeLogger struct {
	ctx context.Context
}

// NewAELogger returns appengine logger
func NewAELogger(ctx context.Context) core.Logger {
	return &aeLogger{ctx}
}

func (l *aeLogger) Debug(args ...interface{}) {
	l.Debugf("%v", fmt.Sprint(args...))
}

func (l *aeLogger) Debugf(format string, args ...interface{}) {
	log.Debugf(l.ctx, format, args...)
}

func (l *aeLogger) Info(args ...interface{}) {
	l.Infof("%v", fmt.Sprint(args...))
}

func (l *aeLogger) Infof(format string, args ...interface{}) {
	log.Infof(l.ctx, format, args...)
}

func (l *aeLogger) Warning(args ...interface{}) {
	l.Warningf("%v", fmt.Sprint(args...))
}

func (l *aeLogger) Warningf(format string, args ...interface{}) {
	log.Warningf(l.ctx, format, args...)
}

func (l *aeLogger) Error(args ...interface{}) {
	l.Errorf("%v", fmt.Sprint(args...))
}

func (l *aeLogger) Errorf(format string, args ...interface{}) {
	log.Errorf(l.ctx, format, args...)
}

func (l *aeLogger) Critical(args ...interface{}) {
	l.Criticalf("%v", fmt.Sprint(args...))
}

func (l *aeLogger) Criticalf(format string, args ...interface{}) {
	log.Criticalf(l.ctx, format, args...)
}
