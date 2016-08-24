package app

import (
	"fmt"
	"log"
	"sync"

	"golang.org/x/net/context"
	glog "google.golang.org/appengine/log"
	mlog "github.com/utahta/momoclo-channel/log"
)

type gaeLogger struct {
	context context.Context
}

var gaeLog mlog.Logger

func NewGaeLogger(ctx context.Context) mlog.Logger {
	return &gaeLogger{context: ctx}
}
func GaeLog(ctx context.Context) mlog.Logger {
	m := new(sync.Mutex)
	m.Lock()
	defer m.Unlock()
	if gaeLog != nil {
		return gaeLog
	}
	gaeLog = NewGaeLogger(ctx)
	return gaeLog
}

func (l gaeLogger) Fatal(args ...interface{}) {
	log.Fatal(fmt.Sprintf("FATAL: %s", args...))
}
func (l gaeLogger) Fatalf(format string, args ...interface{}) {
	log.Fatalf("FATAL: %s", fmt.Sprintf(format, args...))
}
func (l gaeLogger) Panic(args ...interface{}) {
	log.Panic(fmt.Sprintf("PANIC: %s", args...))
}
func (l gaeLogger) Panicf(format string, args ...interface{}) {
	log.Panicf("PANIC: %s", fmt.Sprintf(format, args...))
}
func (l gaeLogger) Critical(args ...interface{}) {
	glog.Criticalf(l.context, "%s", fmt.Sprint(args...))
}
func (l gaeLogger) Criticalf(format string, args ...interface{}) {
	glog.Criticalf(l.context, format, args...)
}
func (l gaeLogger) Error(args ...interface{}) {
	glog.Errorf(l.context, "%s", fmt.Sprint(args...))
}
func (l gaeLogger) Errorf(format string, args ...interface{}) {
	glog.Errorf(l.context, format, args...)
}
func (l gaeLogger) Warning(args ...interface{}) {
	glog.Warningf(l.context, "%s", fmt.Sprint(args...))
}
func (l gaeLogger) Warningf(format string, args ...interface{}) {
	glog.Warningf(l.context, format, args...)
}
func (l gaeLogger) Notice(args ...interface{}) {
	glog.Warningf(l.context, "%s", fmt.Sprint(args...)) // instead of warning
}
func (l gaeLogger) Noticef(format string, args ...interface{}) {
	glog.Warningf(l.context, format, args...) // instead of warning
}
func (l gaeLogger) Info(args ...interface{}) {
	glog.Infof(l.context, "%s", fmt.Sprint(args...))
}
func (l gaeLogger) Infof(format string, args ...interface{}) {
	glog.Infof(l.context, format, args...)
}
func (l gaeLogger) Debug(args ...interface{}) {
	glog.Debugf(l.context, "%s", fmt.Sprint(args...))
}
func (l gaeLogger) Debugf(format string, args ...interface{}) {
	glog.Debugf(l.context, format, args...)
}
