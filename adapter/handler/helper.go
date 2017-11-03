package handler

import (
	"context"
	"net/http"

	"google.golang.org/appengine/log"
)

func Fail(ctx context.Context, w http.ResponseWriter, err error, code int) {
	var message string
	if err != nil {
		message = err.Error()
	}

	log.Errorf(ctx, "An error has occurred! code:%v err:%+v", code, err)
	http.Error(w, message, code)
}
