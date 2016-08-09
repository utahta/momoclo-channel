package momoclo_channel

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func tweetHandler(w http.ResponseWriter, r *http.Request) *appError {
	ctx := appengine.NewContext(r)

	log.Infof(ctx, r.FormValue("items"))
	return nil
}
