package linenotify

import (
	"bytes"
	"fmt"
	"io"
	"sync/atomic"

	"github.com/utahta/go-linenotify"
	"github.com/utahta/momoclo-channel/lib/backoff"
	"github.com/utahta/momoclo-channel/lib/config"
	"github.com/utahta/momoclo-channel/lib/log"
	"github.com/utahta/momoclo-channel/model"
	"github.com/utahta/momoclo-channel/model/linenotification"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
	"google.golang.org/appengine/urlfetch"
)

type client struct {
	*linenotify.Client
	users   []*model.LineNotification
	context context.Context
}

func newClient(ctx context.Context) (*client, error) {
	c := &client{
		Client:  linenotify.New(),
		context: ctx,
	}
	c.HTTPClient = urlfetch.Client(ctx)

	users, err := linenotification.Repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	c.users = users

	return c, nil
}

func (c *client) notifyChannelItem(param *ChannelParam) error {
	message := fmt.Sprintf("\n%s\n%s\n%s", param.Title, param.Item.Title, param.Item.Url)

	if len(param.Item.Images) > 0 {
		image := param.Item.Images[0]
		if err := c.notifyMessage(message, image.Url); err != nil {
			return err
		}

		for _, image := range param.Item.Images[1:] {
			c.notifyMessage(" ", image.Url) // go on
		}
	} else {
		return c.notifyMessage(message, "")
	}
	return nil
}

func (c *client) notifyMessage(message, imageURL string) error {
	if config.C.Linenotify.Disabled {
		return nil
	}

	// prepare cache image
	if imageURL != "" {
		_, err := fetchImage(c.HTTPClient, imageURL)
		if err != nil {
			return err
		}
		defer clearImage(imageURL)
	}

	var (
		ctx             = c.context
		count     int32 = 0
		workQueue       = make(chan bool, 30) // max goroutine
	)
	defer close(workQueue)

	eg := &errgroup.Group{}
	for _, user := range c.users {
		workQueue <- true
		user := user

		eg.Go(func() error {
			defer func() {
				<-workQueue
			}()

			token, err := user.Token()
			if err != nil {
				log.Errorf(ctx, "Failed to get token. hash:%v err:%v", user.Id, err)
				return err
			}

			var image io.Reader
			if b := cacheImage(imageURL); b != nil {
				image = bytes.NewReader(b)
			}

			err = backoff.Retry(3, func() error {
				err := c.Notify(token, message, "", "", image)
				if err == linenotify.ErrNotifyInvalidAccessToken {
					err = nil
					user.Delete(c.context)
					log.Infof(ctx, "Delete LINE Notify token. user:%#v", user)
				}
				return err
			})
			if err != nil {
				log.Errorf(ctx, "Failed to LINE Notify. hash:%v err:%v", user.Id, err)
				return err
			}
			atomic.AddInt32(&count, 1)
			return nil
		})
	}
	eg.Wait()

	log.Infof(ctx, "LINE Notify. message:%s imageURL:%s len:%v/%d", message, imageURL, count, len(c.users))
	return nil
}
