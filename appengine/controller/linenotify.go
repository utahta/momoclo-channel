package controller

import (
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/utahta/go-linenotify"
	"github.com/utahta/momoclo-channel/appengine/lib/log"
	"github.com/utahta/momoclo-channel/appengine/model"
	"google.golang.org/appengine/urlfetch"
)

// LINE Notify と連携する
func LinenotifyOn(w http.ResponseWriter, req *http.Request) {
	ctx := getContext(req)

	c, err := linenotify.NewAuthorization(os.Getenv("LINENOTIFY_CLIENT_ID"), buildURL(req.URL, "/linenotify/callback"))
	if err != nil {
		newError(err, http.StatusInternalServerError).Handle(ctx, w)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "state", Value: c.State, Expires: time.Now().Add(60 * time.Second), Secure: true})

	err = c.Redirect(w, req)
	if err != nil {
		newError(err, http.StatusInternalServerError).Handle(ctx, w)
		return
	}
}

// LINE Notify の連携を解除する
func LinenotifyOff(w http.ResponseWriter, req *http.Request) {
	ctx := getContext(req)
	log.GaeLog(ctx).Info("Redirect to LINE Notification revoke page")

	// official url
	http.Redirect(w, req, "https://notify-bot.line.me/my/", http.StatusFound)
}

func LinenotifyCallback(w http.ResponseWriter, req *http.Request) {
	ctx := getContext(req)

	params, err := linenotify.ParseAuthorization(req)
	if err != nil {
		newError(err, http.StatusInternalServerError).Handle(ctx, w)
		return
	}

	state, err := req.Cookie("state")
	if err != nil {
		newError(err, http.StatusInternalServerError).Handle(ctx, w)
		return
	}

	if params.State != state.Value {
		newError(errors.New("Invalid csrf token."), http.StatusBadRequest).Handle(ctx, w)
		return
	}

	c := linenotify.NewToken(
		params.Code,
		buildURL(req.URL, "/linenotify/callback"),
		os.Getenv("LINENOTIFY_CLIENT_ID"),
		os.Getenv("LINENOTIFY_CLIENT_SECRET"),
	)
	c.HTTPClient = urlfetch.Client(ctx)

	token, err := c.Get()
	if err != nil {
		newError(err, http.StatusInternalServerError).Handle(ctx, w)
		return
	}

	ln, err := model.NewLineNotification(token)
	if err != nil {
		newError(err, http.StatusInternalServerError).Handle(ctx, w)
		return
	}
	ln.Put(ctx) // save to datastore

	t, err := template.New("callback").Parse("<html><body><h1>通知ノフ設定オンにしました（・Θ・）</h1></body></html>")
	if err != nil {
		newError(err, http.StatusInternalServerError).Handle(ctx, w)
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		newError(err, http.StatusInternalServerError).Handle(ctx, w)
		return
	}

	log.GaeLog(ctx).Infof("LINE Notification accepted! id:%v", ln.Id)
}
