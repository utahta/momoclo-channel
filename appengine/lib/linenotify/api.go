package linenotify

import (
	"fmt"
	"sync"

	"github.com/utahta/momoclo-channel/appengine/lib/log"
	"github.com/utahta/momoclo-channel/appengine/model"
	"github.com/utahta/momoclo-channel/linenotify"
	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
)

func notifyMessage(ctx context.Context, message, imageThumbnail, imageFullsize string) {
	glog := log.GaeLog(ctx)
	if message != "" {
		message = fmt.Sprintf("\n%s", message) // [Notify Name] が先頭に入るので改行して調整
	}

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
				glog.Error(err)
			}
		}(item)
	}
	wg.Wait()
}

func NotifyMessage(ctx context.Context, message string) {
	notifyMessage(ctx, message, "", "")
}
