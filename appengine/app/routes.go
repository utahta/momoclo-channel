package app

import (
	"net/http"
	"time"

	"github.com/utahta/momoclo-channel/appengine/controller"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
)

func initRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(appengine.NewContext(req), 55*time.Second)
		defer cancel()
		var err *controller.Error

		switch req.URL.Path {
		case "/cron/crawl":
			err = controller.CronCrawl(ctx, w, req)
		case "/cron/ustream":
			err = controller.CronUstream(ctx, w, req)
		case "/cron/reminder":
			err = controller.CronReminder(ctx, w, req)

		case "/queue/tweet":
			err = controller.QueueTweet(ctx, w, req)
		case "/queue/line":
			err = controller.QueueLine(ctx, w, req)

		case "/linebot/callback":
			err = controller.LineBotCallback(ctx, w, req)
		case "/linebot/help":
			err = controller.LineBotHelp(ctx, w, req)
		case "/linebot/about":
			err = controller.LineBotAbout(ctx, w, req)

		case "/linenotify/on":
			err = controller.LinenotifyOn(ctx, w, req)
		case "/linenotify/off":
			err = controller.LinenotifyOff(ctx, w, req)
		case "/linenotify/callback":
			err = controller.LinenotifyCallback(ctx, w, req)

		default:
			http.NotFound(w, req)
		}
		err.Handle(ctx, w)
	})
}
