package controller

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	"github.com/utahta/momoclo-channel/appengine/lib/log"
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
	log.GaeLog(ctx).Errorf("error:%+v code:%d", e.Error, e.Code)
}

func buildURL(u *url.URL, path string) string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%s://%s", u.Scheme, u.Host))
	if len(path) > 0 && path[0] != '/' {
		buf.WriteString("/")
	}
	buf.WriteString(path)

	return buf.String()
}
