package momoclo_channel

import (
	"net/http"
	"encoding/json"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"github.com/utahta/momoclo-channel/crawler"
	"github.com/pkg/errors"
)

func tweetHandler(w http.ResponseWriter, r *http.Request) *appError {
	ctx := appengine.NewContext(r)

	items := []*crawler.ChannelItem{}
	if err := json.Unmarshal([]byte(r.FormValue("items")), &items); err != nil {
		return newAppError(errors.Wrap(err, "Failed to unmarshal items."), http.StatusInternalServerError)
	}

	for _, item := range items {
		log.Infof(ctx, "%v", item)
	}
	return nil
}
