package app

import (
	"net/http"
	"time"

	"github.com/utahta/momoclo-channel/appengine/controller"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
)

func initRoutes() {
	http.Handle("/cron/", &CronHandler{})
	http.Handle("/queue/", &QueueHandler{})
	http.Handle("/linenotify/", &LinenotifyHandler{})

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(appengine.NewContext(req), 55*time.Second)
		defer cancel()
		var err *controller.Error

		switch req.URL.Path {
		case "/linebot/callback":
			err = controller.LineBotCallback(ctx, w, req)
		case "/linebot/help":
			err = controller.LineBotHelp(ctx, w, req)
		default:
			http.NotFound(w, req)
		}
		err.Handle(ctx, w)
	})
}
