package linenotify

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/appengine/lib/log"
	"github.com/utahta/momoclo-channel/appengine/lib/util"
	"github.com/utahta/momoclo-channel/appengine/model"
	"github.com/utahta/momoclo-channel/crawler"
	"github.com/utahta/momoclo-channel/linenotify"
	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
)

func notifyMessage(ctx context.Context, message, imageThumbnail, imageFullsize string) error {
	glog := log.NewGaeLogger(ctx)

	query := model.NewLineNotificationQuery(ctx)
	items, err := query.GetAll()
	if err != nil {
		return err
	}

	var errCount int32
	req := linenotify.NewRequestNotify()
	req.Client = urlfetch.Client(ctx)

	var workQueue = make(chan bool, 10) // max goroutine
	var wg sync.WaitGroup
	for _, item := range items {
		workQueue <- true
		wg.Add(1)
		go func(item *model.LineNotification) {
			defer wg.Done()
			defer func() {
				<-workQueue
			}()

			token, err := item.Token()
			if err != nil {
				glog.Error(err)
				atomic.AddInt32(&errCount, 1)
				return
			}
			if err := req.Notify(token, message, imageThumbnail, imageFullsize); err != nil {
				if err == linenotify.ErrorNotifyInvalidAccessToken {
					item.Delete(ctx)
					glog.Infof("Delete LINE Notify token. hash:%s", item.Id)
				} else {
					glog.Error(err)
				}
				atomic.AddInt32(&errCount, 1)
			}
		}(item)
	}
	wg.Wait()

	glog.Infof("LINE Notify. message:%s imageURL:%s len:%d errCount:%d", message, imageFullsize, len(items), errCount)
	return nil
}

// Send message to LINE Notify
func NotifyMessage(ctx context.Context, message string) error {
	// [Notify Name] が付くので先頭に改行をいれて調整
	return notifyMessage(ctx, fmt.Sprintf("\n%s", message), "", "")
}

// Send channel message and images to LINE Notify
func NotifyChannel(ctx context.Context, ch *crawler.Channel) error {
	errFlg := util.NewAtomicBool(false)
	var wg sync.WaitGroup
	wg.Add(len(ch.Items))
	for _, item := range ch.Items {
		go func(ctx context.Context, item *crawler.ChannelItem) {
			defer wg.Done()

			if err := model.NewLineItem(item).Put(ctx); err != nil {
				return
			}
			if err := notifyChannelItem(ctx, ch.Title, item); err != nil {
				errFlg.Set(true)
				log.GaeLog(ctx).Error(err)
			}
		}(ctx, item)
	}
	wg.Wait()

	if errFlg.Enabled() {
		return errors.New("Errors occured in linenotify.NotifyChannel")
	}
	return nil
}

func notifyChannelItem(ctx context.Context, title string, item *crawler.ChannelItem) error {
	message := fmt.Sprintf("\n%s\n%s\n%s", title, item.Title, item.Url)

	if len(item.Images) > 0 {
		image := item.Images[0]
		if err := notifyMessage(ctx, message, image.Url, image.Url); err != nil {
			return err
		}

		for _, image := range item.Images[1:] {
			if err := notifyMessage(ctx, " ", image.Url, image.Url); err != nil {
				return err
			}
		}
	} else {
		return notifyMessage(ctx, message, "", "")
	}
	return nil
}
