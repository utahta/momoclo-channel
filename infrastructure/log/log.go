package log

import (
	"context"
	"fmt"

	"github.com/utahta/momoclo-channel/domain/event"
	"google.golang.org/appengine/log"
)

type appengineLogger struct {
	ctx context.Context
}

// NewAppengineLogger returns appengine logger
func NewAppengineLogger(ctx context.Context) event.Logger {
	return &appengineLogger{ctx}
}

func (l *appengineLogger) Debug(args ...interface{}) {
	l.Debugf("%v", fmt.Sprint(args...))
}

func (l *appengineLogger) Debugf(format string, args ...interface{}) {
	log.Debugf(l.ctx, format, args...)
}

func (l *appengineLogger) Info(args ...interface{}) {
	l.Infof("%v", fmt.Sprint(args...))
}

func (l *appengineLogger) Infof(format string, args ...interface{}) {
	log.Infof(l.ctx, format, args...)
}

func (l *appengineLogger) Warning(args ...interface{}) {
	l.Warningf("%v", fmt.Sprint(args...))
}

func (l *appengineLogger) Warningf(format string, args ...interface{}) {
	log.Warningf(l.ctx, format, args...)
}

func (l *appengineLogger) Error(args ...interface{}) {
	l.Errorf("%v", fmt.Sprint(args...))
}

func (l *appengineLogger) Errorf(format string, args ...interface{}) {
	log.Errorf(l.ctx, format, args...)
}

func (l *appengineLogger) Critical(args ...interface{}) {
	l.Criticalf("%v", fmt.Sprint(args...))
}

func (l *appengineLogger) Criticalf(format string, args ...interface{}) {
	log.Criticalf(l.ctx, format, args...)
}
