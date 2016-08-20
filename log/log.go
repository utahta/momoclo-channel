package log

import (
	"log"
	"fmt"

	glog "google.golang.org/appengine/log"
	"golang.org/x/net/context"
)

type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Criticalf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
}

type silentLogger struct {
}

func NewSilentLogger() Logger {
	return &silentLogger{}
}

func (_ silentLogger) Debugf(_ string, _ ...interface{}) {}
func (_ silentLogger) Infof(_ string, _ ...interface{}) {}
func (_ silentLogger) Warningf(_ string, _ ...interface{}) {}
func (_ silentLogger) Errorf(_ string, _ ...interface{}) {}
func (_ silentLogger) Criticalf(_ string, _ ...interface{}) {}
func (_ silentLogger) Fatalf(_ string, _ ...interface{}) {}
func (_ silentLogger) Panicf(_ string, _ ...interface{}) {}

type basicLogger struct {
}
var basicLog Logger

func NewBasicLogger() Logger {
	return &basicLogger{}
}
func Basic() Logger {
	if basicLog != nil {
		return basicLog
	}
	basicLog = NewBasicLogger()
	return basicLog
}

func (_ basicLogger) Debugf(format string, args ...interface{}) { log.Printf("DEBUG: %s", fmt.Sprintf(format, args...)) }
func (_ basicLogger) Infof(format string, args ...interface{}) { log.Printf("INFO: %s", fmt.Sprintf(format, args...)) }
func (_ basicLogger) Warningf(format string, args ...interface{}) { log.Printf("WARN: %s", fmt.Sprintf(format, args...)) }
func (_ basicLogger) Errorf(format string, args ...interface{}) { log.Printf("ERROR: %s", fmt.Sprintf(format, args...)) }
func (_ basicLogger) Criticalf(format string, args ...interface{}) { log.Printf("CRIT: %s", fmt.Sprintf(format, args...)) }
func (_ basicLogger) Fatalf(format string, args ...interface{}) { log.Fatalf("FATAL: %s", fmt.Sprintf(format, args...)) }
func (_ basicLogger) Panicf(format string, args ...interface{}) { log.Panicf("PANIC: %s", fmt.Sprintf(format, args...)) }

type gaeLogger struct {
	context context.Context
}
var gaeLog Logger

func NewGaeLogger(ctx context.Context) Logger {
	return &gaeLogger{ context: ctx }
}
func Gae(ctx context.Context) Logger {
	if gaeLog != nil {
		return gaeLog
	}
	gaeLog = NewGaeLogger(ctx)
	return gaeLog
}

func (l gaeLogger) Debugf(format string, args ...interface{}) { glog.Debugf(l.context, format, args...) }
func (l gaeLogger) Infof(format string, args ...interface{}) { glog.Infof(l.context, format, args...) }
func (l gaeLogger) Warningf(format string, args ...interface{}) { glog.Warningf(l.context, format, args...) }
func (l gaeLogger) Errorf(format string, args ...interface{}) { glog.Errorf(l.context, format, args...) }
func (l gaeLogger) Criticalf(format string, args ...interface{}) { glog.Criticalf(l.context, format, args...) }
func (l gaeLogger) Fatalf(format string, args ...interface{}) { log.Fatalf("FATAL: %s", fmt.Sprintf(format, args...)) }
func (l gaeLogger) Panicf(format string, args ...interface{}) { log.Panicf("PANIC: %s", fmt.Sprintf(format, args...)) }
