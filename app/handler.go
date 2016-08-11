package momoclo_channel

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

type appError struct {
	Error error
	Code int
}

func newAppError(err error, code int) *appError {
	return &appError{ Error: err, Code: code }
}

type appHandler func(w http.ResponseWriter, r *http.Request) *appError

func appHandlerFunc(fn appHandler) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		err := fn(w, r)
		if err != nil {
			http.Error(w, err.Error.Error(), err.Code)
			log.Errorf(appengine.NewContext(r), "error:%v code:%d", err.Error, err.Code)
		}
	}
}
