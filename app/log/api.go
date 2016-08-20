package log

import (
	"log"
	"fmt"

	"golang.org/x/net/context"
	gaelog "google.golang.org/appengine/log"
)

func Debugf(ctx context.Context, format string, args ...interface{}) {
	if ctx != nil {
		gaelog.Debugf(ctx, format, args...)
	} else {
		log.Printf("DEBUG: %s", fmt.Sprintf(format, args...))
	}
}

// Infof is like Debugf, but at Info level.
func Infof(ctx context.Context, format string, args ...interface{}) {
	if ctx != nil {
		gaelog.Infof(ctx, format, args...)
	} else {
		log.Printf("INFO: %s", fmt.Sprintf(format, args...))
	}
}

// Warningf is like Debugf, but at Warning level.
func Warningf(ctx context.Context, format string, args ...interface{}) {
	if ctx != nil {
		gaelog.Warningf(ctx, format, args...)
	} else {
		log.Printf("WARN: %s", fmt.Sprintf(format, args...))
	}
}

// Errorf is like Debugf, but at Error level.
func Errorf(ctx context.Context, format string, args ...interface{}) {
	if ctx != nil {
		gaelog.Errorf(ctx, format, args...)
	} else {
		log.Printf("ERROR: %s", fmt.Sprintf(format, args...))
	}
}

// Criticalf is like Debugf, but at Critical level.
func Criticalf(ctx context.Context, format string, args ...interface{}) {
	if ctx != nil {
		gaelog.Criticalf(ctx, format, args...)
	} else {
		log.Printf("CRIT: %s", fmt.Sprintf(format, args...))
	}
}
