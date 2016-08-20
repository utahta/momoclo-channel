package app

import (
	"net/http"

	"github.com/utahta/momoclo-channel/log"
	"golang.org/x/net/context"
)

type Error struct {
	Error error
	Code  int
}

func newError(err error, code int) *Error {
	return &Error{Error: err, Code: code}
}

func (e *Error) Handle(ctx context.Context, w http.ResponseWriter) {
	if e == nil {
		return
	}
	http.Error(w, e.Error.Error(), e.Code)
	log.Gae(ctx).Errorf("error:%v code:%d", e.Error, e.Code)
}
