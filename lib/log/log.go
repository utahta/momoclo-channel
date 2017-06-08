package log

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

func Debug(ctx context.Context, args ...interface{}) {
	Debugf(ctx, "%v", fmt.Sprint(args...))
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	log.Debugf(ctx, format, args...)
}

func Info(ctx context.Context, args ...interface{}) {
	Infof(ctx, "%v", fmt.Sprint(args...))
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	log.Infof(ctx, format, args...)
}

func Warning(ctx context.Context, args ...interface{}) {
	Warningf(ctx, "%v", fmt.Sprint(args...))
}

func Warningf(ctx context.Context, format string, args ...interface{}) {
	log.Warningf(ctx, format, args...)
}

func Error(ctx context.Context, args ...interface{}) {
	Errorf(ctx, "%v", fmt.Sprint(args...))
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	log.Errorf(ctx, format, args...)
}

func Critical(ctx context.Context, args ...interface{}) {
	Criticalf(ctx, "%v", fmt.Sprint(args...))
}

func Criticalf(ctx context.Context, format string, args ...interface{}) {
	log.Criticalf(ctx, format, args...)
}
