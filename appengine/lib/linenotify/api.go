package linenotify

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/utahta/go-atomicbool"
	"github.com/utahta/momoclo-channel/appengine/lib/log"
	"github.com/utahta/momoclo-channel/appengine/model"
	"github.com/utahta/momoclo-channel/linenotify"
	"github.com/utahta/momoclo-crawler"
	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
)

// Send message to LINE Notify
func NotifyMessage(ctx context.Context, message string) error {
	if disabled() {
		return nil
	}

	// [Notify Name] が付くので先頭に改行をいれて調整
	return notifyMessage(ctx, fmt.Sprintf("\n%s", message), "")
}

// Send channel message and images to LINE Notify
func NotifyChannel(ctx context.Context, ch *crawler.Channel) error {
	if disabled() {
		return nil
	}

	errFlg := atomicbool.New(false)
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

func notifyMessage(ctx context.Context, message, imageFile string) error {
	glog := log.NewGaeLogger(ctx)

	query := model.NewLineNotificationQuery(ctx)
	items, err := query.GetAll()
	if err != nil {
		return err
	}

	var errCount int32
	reqCtx, cancel := context.WithTimeout(ctx, 540*time.Second)
	defer cancel()
	req := linenotify.NewRequestNotify()
	req.Client = urlfetch.Client(reqCtx)

	// 先にキャッシュしておく
	if imageFile != "" {
		_, err := req.CacheImage(imageFile)
		if err != nil {
			return err
		}
		defer req.ClearImage(imageFile)
	}

	var workQueue = make(chan bool, 1000) // max goroutine
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
			if err := req.Notify(token, message, "", "", imageFile); err != nil {
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

	glog.Infof("LINE Notify. message:%s imageURL:%s len:%d errCount:%d", message, imageFile, len(items), errCount)
	return nil
}

func notifyChannelItem(ctx context.Context, title string, item *crawler.ChannelItem) error {
	message := fmt.Sprintf("\n%s\n%s\n%s", title, item.Title, item.Url)

	if len(item.Images) > 0 {
		image := item.Images[0]
		if err := notifyMessage(ctx, message, image.Url); err != nil {
			return err
		}

		for _, image := range item.Images[1:] {
			if err := notifyMessage(ctx, " ", image.Url); err != nil {
				return err
			}
		}
	} else {
		return notifyMessage(ctx, message, "")
	}
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
