package ustream

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/appengine/lib/linenotify"
	"github.com/utahta/momoclo-channel/appengine/lib/twitter"
	"github.com/utahta/momoclo-channel/appengine/model"
	"github.com/utahta/momoclo-channel/ustream"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/urlfetch"
)

func Notify(ctx context.Context) error {
	c := ustream.NewClient()
	c.HttpClient.Transport = &urlfetch.Transport{Context: ctx}

	isLive, err := c.IsLive()
	if err != nil {
		return errors.Wrap(err, "Failed to get ustream status")
	}

	status := model.NewUstreamStatus()
	if err := status.Get(ctx); err != nil && err != datastore.ErrNoSuchEntity {
		return errors.Wrap(err, "Failed to get ustream status from datastore")
	}

	if status.IsLive == isLive {
		return nil
	}
	status.IsLive = isLive
	status.Put(ctx)

	if isLive {
		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			twitter.TweetUstream(ctx)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			linenotify.NotifyMessage(ctx, "momocloTV が配信を開始しました")
		}()

		wg.Wait()
	}
	return nil
}
