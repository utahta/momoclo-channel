package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
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

func (l silentLogger) Fatal(_ ...interface{})               {}
func (l silentLogger) Fatalf(_ string, _ ...interface{})    {}
func (l silentLogger) Panic(_ ...interface{})               {}
func (l silentLogger) Panicf(_ string, _ ...interface{})    {}
func (l silentLogger) Critical(_ ...interface{})            {}
func (l silentLogger) Criticalf(_ string, _ ...interface{}) {}
func (l silentLogger) Error(_ ...interface{})               {}
func (l silentLogger) Errorf(_ string, _ ...interface{})    {}
func (l silentLogger) Warning(_ ...interface{})             {}
func (l silentLogger) Warningf(_ string, _ ...interface{})  {}
func (l silentLogger) Notice(_ ...interface{})              {}
func (l silentLogger) Noticef(_ string, _ ...interface{})   {}
func (l silentLogger) Info(_ ...interface{})                {}
func (l silentLogger) Infof(_ string, _ ...interface{})     {}
func (l silentLogger) Debug(_ ...interface{})               {}
func (l silentLogger) Debugf(_ string, _ ...interface{})    {}

type logger struct {
	log *log.Logger
}

func NewIOLogger(io io.Writer) Logger {
	return &logger{log: log.New(io, "", log.LstdFlags)}
}

var basicLog Logger

func NewBasicLogger() Logger {
	return NewIOLogger(os.Stderr)
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

func (l logger) Fatal(args ...interface{}) {
	l.log.Fatal(fmt.Sprintf("FATAL: %+v", args...))
}
func (l logger) Fatalf(format string, args ...interface{}) {
	l.log.Fatalf("FATAL: %s", fmt.Sprintf(format, args...))
}
func (l logger) Panic(args ...interface{}) {
	l.log.Panic(fmt.Sprintf("PANIC: %+v", args...))
}
func (l logger) Panicf(format string, args ...interface{}) {
	l.log.Panicf("PANIC: %s", fmt.Sprintf(format, args...))
}
func (l logger) Critical(args ...interface{}) {
	l.log.Print(fmt.Sprintf("CRIT: %+v", args...))
}
func (l logger) Criticalf(format string, args ...interface{}) {
	l.log.Printf("CRIT: %s", fmt.Sprintf(format, args...))
}
func (l logger) Error(args ...interface{}) {
	l.log.Print(fmt.Sprintf("ERROR: %+v", args...))
}
func (l logger) Errorf(format string, args ...interface{}) {
	l.log.Printf("ERROR: %s", fmt.Sprintf(format, args...))
}
func (l logger) Warning(args ...interface{}) {
	l.log.Print(fmt.Sprintf("WARN: %s", args...))
}
func (l logger) Warningf(format string, args ...interface{}) {
	l.log.Printf("WARN: %s", fmt.Sprintf(format, args...))
}
func (l logger) Notice(args ...interface{}) {
	l.log.Print(fmt.Sprintf("NOTICE: %s", args...))
}
func (l logger) Noticef(format string, args ...interface{}) {
	l.log.Printf("NOTICE: %s", fmt.Sprintf(format, args...))
}
func (l logger) Info(args ...interface{}) {
	l.log.Print(fmt.Sprintf("INFO: %s", args...))
}
func (l logger) Infof(format string, args ...interface{}) {
	l.log.Printf("INFO: %s", fmt.Sprintf(format, args...))
}
func (l logger) Debug(args ...interface{}) {
	l.log.Print(fmt.Sprintf("DEBUG: %s", args...))
}
func (l logger) Debugf(format string, args ...interface{}) {
	l.log.Printf("DEBUG: %s", fmt.Sprintf(format, args...))
}
