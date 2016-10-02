package linenotify

import (
	"fmt"
	"sync"

	"github.com/utahta/momoclo-channel/appengine/lib/log"
	"github.com/utahta/momoclo-channel/appengine/model"
	"github.com/utahta/momoclo-channel/crawler"
	"github.com/utahta/momoclo-channel/linenotify"
	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
)

func notifyMessage(ctx context.Context, message, imageThumbnail, imageFullsize string) {
	glog := log.GaeLog(ctx)

	query := model.NewLineNotificationQuery(ctx)
	items, err := query.GetAll()
	if err != nil {
		glog.Error(err)
		return
	}

	req := linenotify.NewRequestNotify()
	req.Client = urlfetch.Client(ctx)

	var workQueue = make(chan bool, 20) // max goroutine
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
				return
			}
			if err := req.Notify(token, message, imageThumbnail, imageFullsize); err != nil {
				if err == linenotify.ErrorNotifyInvalidAccessToken {
					item.Delete(ctx)
					glog.Infof("Delete LINE Notify token. hash:%s", item.Id)
				} else {
					glog.Error(err)
					return
				}
			}
		}(item)
	}
	wg.Wait()

	glog.Infof("LINE Notify. message:%s imageURL:%s len:%d", message, imageFullsize, len(items))
}

func NotifyMessage(ctx context.Context, message string) {
	// [Notify Name] が付くので先頭に改行をいれて調整
	notifyMessage(ctx, fmt.Sprintf("\n%s", message), "", "")
}

func NotifyChannelItem(ctx context.Context, title string, item *crawler.ChannelItem) {
	message := fmt.Sprintf("\n%s\n%s\n%s", title, item.Title, item.Url)

	if len(item.Images) > 0 {
		image := item.Images[0]
		notifyMessage(ctx, message, image.Url, image.Url)

		for _, image := range item.Images[1:] {
			notifyMessage(ctx, " ", image.Url, image.Url)
		}
	} else {
		notifyMessage(ctx, message, "", "")
	}
}
