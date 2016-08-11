package app

import (
	"net/http"

	"google.golang.org/appengine/log"
	"golang.org/x/net/context"
)

type Error struct {
	Error error
	Code int
}

func newError(err error, code int) *Error {
	return &Error{ Error: err, Code: code }
}

func (e *Error) Handle(ctx context.Context, w http.ResponseWriter) {
	if e == nil {
		return
	}
	http.Error(w, e.Error.Error(), e.Code)
	log.Errorf(ctx, "error:%v code:%d", e.Error, e.Code)
}
