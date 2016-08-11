package momoclo_channel

import (
	"net/http"

	"google.golang.org/appengine/log"
	"golang.org/x/net/context"
)

func appError(ctx context.Context, w http.ResponseWriter, err error, code int) {
	http.Error(w, err.Error(), code)
	log.Errorf(ctx, "error:%v code:%d", err, code)
}
