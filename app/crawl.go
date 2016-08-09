package momoclo_channel

import (
	"net/http"
	"net/url"
	"encoding/json"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"google.golang.org/appengine/taskqueue"
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/crawler"
)

func crawlHandler(w http.ResponseWriter, r *http.Request) *appError {
	ctx := appengine.NewContext(r)

	c := crawler.NewMomotaBlogChannel()
	c.HttpClient.Transport = &urlfetch.Transport{Context: ctx}
	items, err := crawler.FetchParse(c)
	if err != nil {
		return newAppError(errors.Wrapf(err, "Failed to fetch momota blog"), http.StatusInternalServerError)
	}

	bin, err := json.Marshal(items)
	if err != nil {
		return newAppError(errors.Wrapf(err, "Failed to json encode from momota blog"), http.StatusInternalServerError)
	}

	task := taskqueue.NewPOSTTask("/queue/tweet", url.Values{
		"items": {string(bin)},
	})
	_, err = taskqueue.Add(ctx, task, "queue-tweet")
	if err != nil {
		return newAppError(errors.Wrapf(err, "Failed to taskqueue.Add"), http.StatusInternalServerError)
	}
	return nil
}
