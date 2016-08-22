package log

import (
	"fmt"
	"log"
	"sync"

	"golang.org/x/net/context"
	glog "google.golang.org/appengine/log"
)

type Logger interface {
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})

	Panic(args ...interface{})
	Panicf(format string, args ...interface{})

	Critical(args ...interface{})
	Criticalf(format string, args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	Warning(args ...interface{})
	Warningf(format string, args ...interface{})

	Notice(args ...interface{})
	Noticef(format string, args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})

	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
}

type silentLogger struct {
}

func NewSilentLogger() Logger {
	return &silentLogger{}
}

func (_ silentLogger) Fatal(_ ...interface{})               {}
func (_ silentLogger) Fatalf(_ string, _ ...interface{})    {}
func (_ silentLogger) Panic(_ ...interface{})               {}
func (_ silentLogger) Panicf(_ string, _ ...interface{})    {}
func (_ silentLogger) Critical(_ ...interface{})            {}
func (_ silentLogger) Criticalf(_ string, _ ...interface{}) {}
func (_ silentLogger) Error(_ ...interface{})               {}
func (_ silentLogger) Errorf(_ string, _ ...interface{})    {}
func (_ silentLogger) Warning(_ ...interface{})             {}
func (_ silentLogger) Warningf(_ string, _ ...interface{})  {}
func (_ silentLogger) Notice(_ ...interface{})              {}
func (_ silentLogger) Noticef(_ string, _ ...interface{})   {}
func (_ silentLogger) Info(_ ...interface{})                {}
func (_ silentLogger) Infof(_ string, _ ...interface{})     {}
func (_ silentLogger) Debug(_ ...interface{})               {}
func (_ silentLogger) Debugf(_ string, _ ...interface{})    {}

type basicLogger struct {
}

var basicLog Logger

func NewBasicLogger() Logger {
	return &basicLogger{}
}
func Basic() Logger {
	m := new(sync.Mutex)
	m.Lock()
	defer m.Unlock()
	if basicLog != nil {
		return basicLog
	}
	basicLog = NewBasicLogger()
	return basicLog
}

func (_ basicLogger) Fatal(args ...interface{}) {
	log.Fatal(fmt.Sprintf("FATAL: %s", args...))
}
func (_ basicLogger) Fatalf(format string, args ...interface{}) {
	log.Fatalf("FATAL: %s", fmt.Sprintf(format, args...))
}
func (_ basicLogger) Panic(args ...interface{}) {
	log.Panic(fmt.Sprintf("PANIC: %s", args...))
}
func (_ basicLogger) Panicf(format string, args ...interface{}) {
	log.Panicf("PANIC: %s", fmt.Sprintf(format, args...))
}
func (_ basicLogger) Critical(args ...interface{}) {
	log.Print(fmt.Sprintf("CRIT: %s", args...))
}
func (_ basicLogger) Criticalf(format string, args ...interface{}) {
	log.Printf("CRIT: %s", fmt.Sprintf(format, args...))
}
func (_ basicLogger) Error(args ...interface{}) {
	log.Print(fmt.Sprintf("ERROR: %s", args...))
}
func (_ basicLogger) Errorf(format string, args ...interface{}) {
	log.Printf("ERROR: %s", fmt.Sprintf(format, args...))
}
func (_ basicLogger) Warning(args ...interface{}) {
	log.Print(fmt.Sprintf("WARN: %s", args...))
}
func (_ basicLogger) Warningf(format string, args ...interface{}) {
	log.Printf("WARN: %s", fmt.Sprintf(format, args...))
}
func (_ basicLogger) Notice(args ...interface{}) {
	log.Print(fmt.Sprintf("NOTICE: %s", args...))
}
func (_ basicLogger) Noticef(format string, args ...interface{}) {
	log.Printf("NOTICE: %s", fmt.Sprintf(format, args...))
}
func (_ basicLogger) Info(args ...interface{}) {
	log.Print(fmt.Sprintf("INFO: %s", args...))
}
func (_ basicLogger) Infof(format string, args ...interface{}) {
	log.Printf("INFO: %s", fmt.Sprintf(format, args...))
}
func (_ basicLogger) Debug(args ...interface{}) {
	log.Print(fmt.Sprintf("DEBUG: %s", args...))
}
func (_ basicLogger) Debugf(format string, args ...interface{}) {
	log.Printf("DEBUG: %s", fmt.Sprintf(format, args...))
}

type gaeLogger struct {
	context context.Context
}

var gaeLog Logger

func NewGaeLogger(ctx context.Context) Logger {
	return &gaeLogger{context: ctx}
}
func Gae(ctx context.Context) Logger {
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
