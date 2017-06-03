package linenotify

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/utahta/go-linenotify"
	"github.com/utahta/momoclo-channel/appengine/lib/log"
	"github.com/utahta/momoclo-channel/appengine/model"
	"github.com/utahta/momoclo-crawler"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
	"google.golang.org/appengine/urlfetch"
)

const timeout = 540 * time.Second

// Send message to LINE Notify
func NotifyMessage(ctx context.Context, message string) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	c, err := newClient(ctx)
	if err != nil {
		return err
	}

	// [Notify Name] が付くので先頭に改行をいれて調整
	return c.notifyMessage(fmt.Sprintf("\n%s", message), "")
}

// Send channel message and images to LINE Notify
func NotifyChannel(ctx context.Context, ch *crawler.Channel) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	c, err := newClient(ctx)
	if err != nil {
		return err
	}

	eg := &errgroup.Group{}
	for _, item := range ch.Items {
		item := item
		eg.Go(func() error {
			if err := model.NewLineItem(item).Put(ctx); err != nil {
				return nil
			}
			if err := c.notifyChannelItem(ch.Title, item); err != nil {
				log.GaeLog(ctx).Error(err)
				return err
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return errors.New("Errors occurred in NotifyChannel")
	}
	return nil
}

type client struct {
	*linenotify.Client
	users   []*model.LineNotification
	context context.Context
	log     log.Logger
}

func newClient(ctx context.Context) (*client, error) {
	c := &client{
		Client:  linenotify.New(),
		context: ctx,
		log:     log.NewGaeLogger(ctx),
	}
	c.HTTPClient.Transport = &urlfetch.Transport{Context: ctx}

	query := model.NewLineNotificationQuery(ctx)
	users, err := query.GetAll()
	if err != nil {
		return nil, err
	}
	c.users = users

	return c, nil
}

func (c *client) notifyChannelItem(title string, item *crawler.ChannelItem) error {
	message := fmt.Sprintf("\n%s\n%s\n%s", title, item.Title, item.Url)

	if len(item.Images) > 0 {
		image := item.Images[0]
		if err := c.notifyMessage(message, image.Url); err != nil {
			return err
		}

		for _, image := range item.Images[1:] {
			if err := c.notifyMessage(" ", image.Url); err != nil {
				return err
			}
		}
	} else {
		return c.notifyMessage(message, "")
	}
	return nil
}

func (c *client) notifyMessage(message, imageURL string) error {
	if disabled() {
		return nil
	}

	// prepare cached image
	if imageURL != "" {
		_, err := fetchImage(c.HTTPClient, imageURL)
		if err != nil {
			return err
		}
		defer clearImage(imageURL)
	}

	var workQueue = make(chan bool, 100) // max goroutine
	eg := &errgroup.Group{}
	for _, user := range c.users {
		user := user
		workQueue <- true

		eg.Go(func() error {
			defer func() {
				<-workQueue
			}()

			token, err := user.Token()
			if err != nil {
				c.log.Errorf("Failed to get token. hash:%v err:%v", user.Id, err)
				return err
			}

			var image io.Reader
			if b := cacheImage(imageURL); b != nil {
				image = bytes.NewReader(b)
			}

			err = c.Notify(token, message, "", "", image)
			if err == linenotify.ErrNotifyInvalidAccessToken {
				user.Delete(c.context)
				c.log.Infof("Delete LINE Notify token. hash:%s", user.Id)
			} else if err != nil {
				c.log.Error("Failed to notify. hash:%v err:%v", user.Id, err)
				return err
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		c.log.Warningf("Errors occurred in notifyMessage.")
	}
	c.log.Infof("LINE Notify. message:%s imageURL:%s len:%d", message, imageURL, len(c.users))
	return nil
}

// if true disable linenotify
func disabled() bool {
	e := os.Getenv("LINENOTIFY_DISABLE")
	if e != "" {
		return true
	}
	return false
}
