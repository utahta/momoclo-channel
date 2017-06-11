package controller

import (
	"html/template"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/utahta/go-linenotify"
	"github.com/utahta/momoclo-channel/app"
	"github.com/utahta/momoclo-channel/lib/config"
	"github.com/utahta/momoclo-channel/lib/log"
	"github.com/utahta/momoclo-channel/model"
	"google.golang.org/appengine/urlfetch"
)

// LINE Notify と連携する
func LinenotifyOn(w http.ResponseWriter, req *http.Request) {
	ctx := app.GetContext(req)

	c, err := linenotify.NewAuthorization(config.C.Linenotify.ClientID, config.C.App.BaseURL+"/linenotify/callback")
	if err != nil {
		ctx.Fail(err)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "state", Value: c.State, Expires: time.Now().Add(60 * time.Second), Secure: true})

	err = c.Redirect(w, req)
	if err != nil {
		ctx.Fail(err)
		return
	}
}

// LINE Notify の連携を解除する
func LinenotifyOff(w http.ResponseWriter, req *http.Request) {
	ctx := app.GetContext(req)
	log.Info(ctx, "Redirect to LINE Notification revoke page")

	// official url
	http.Redirect(w, req, "https://notify-bot.line.me/my/", http.StatusFound)
}

func LinenotifyCallback(w http.ResponseWriter, req *http.Request) {
	ctx := app.GetContext(req)

	params, err := linenotify.ParseAuthorization(req)
	if err != nil {
		ctx.Fail(err)
		return
	}

	state, err := req.Cookie("state")
	if err != nil {
		ctx.Fail(err)
		return
	}

	if params.State != state.Value {
		ctx.Error(errors.New("Invalid csrf token."), http.StatusBadRequest)
		return
	}

	c := linenotify.NewToken(
		params.Code,
		config.C.App.BaseURL+"/linenotify/callback",
		config.C.Linenotify.ClientID,
		config.C.Linenotify.ClientSecret,
	)
	c.HTTPClient = urlfetch.Client(ctx)

	token, err := c.Get()
	if err != nil {
		ctx.Fail(err)
		return
	}

	ln, err := model.NewLineNotification(token)
	if err != nil {
		ctx.Fail(err)
		return
	}
	ln.Put(ctx) // save to datastore

	t, err := template.New("callback").Parse("<html><body><h1>通知ノフ設定オンにしました（・Θ・）</h1></body></html>")
	if err != nil {
		ctx.Fail(err)
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		ctx.Fail(err)
		return
	}

	log.Infof(ctx, "LINE Notification accepted! id:%v", ln.Id)
}
