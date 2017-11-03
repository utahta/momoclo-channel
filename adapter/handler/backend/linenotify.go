package backend

import (
	"html/template"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/utahta/go-linenotify"
	"github.com/utahta/momoclo-channel/adapter/handler"
	"github.com/utahta/momoclo-channel/domain/linenotification"
	"github.com/utahta/momoclo-channel/lib/config"
	"github.com/utahta/momoclo-channel/lib/log"
	"google.golang.org/appengine/urlfetch"
)

// LINE Notify と連携する
func LinenotifyOn(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	c, err := linenotify.NewAuthorization(config.C.Linenotify.ClientID, config.C.App.BaseURL+"/linenotify/callback")
	if err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "state", Value: c.State, Expires: time.Now().Add(300 * time.Second), Secure: true})

	err = c.Redirect(w, req)
	if err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}
}

// LINE Notify の連携を解除する
func LinenotifyOff(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	log.Info(ctx, "Redirect to LINE Notification revoke page")

	// official url
	http.Redirect(w, req, "https://notify-bot.line.me/my/", http.StatusFound)
}

func LinenotifyCallback(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	params, err := linenotify.ParseAuthorization(req)
	if err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}

	state, err := req.Cookie("state")
	if err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}

	if params.State != state.Value {
		handler.Fail(ctx, w, errors.New("invalid csrf token"), http.StatusBadRequest)
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
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}

	ln, err := linenotification.Repository.PutToken(ctx, token)
	if err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}

	t, err := template.New("callback").Parse("<html><body><h1>通知ノフ設定オンにしました（・Θ・）</h1></body></html>")
	if err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}

	log.Infof(ctx, "LINE Notification accepted! id:%v", ln.Id)
}
