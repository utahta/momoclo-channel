package controller

import (
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/appengine/model"
	"github.com/utahta/momoclo-channel/linenotify"
	"google.golang.org/appengine/urlfetch"
)

// LINE Notify と連携する
func LinenotifyOn(w http.ResponseWriter, req *http.Request) {
	ctx := getContext(req)

	reqAuth, err := linenotify.NewRequestAuthorization(os.Getenv("LINENOTIFY_CLIENT_ID"), buildURL(req.URL, "/linenotify/callback"))
	if err != nil {
		newError(err, http.StatusInternalServerError).Handle(ctx, w)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "state", Value: reqAuth.State, Expires: time.Now().Add(60 * time.Second), Secure: true})

	err = reqAuth.Redirect(w, req)
	if err != nil {
		newError(err, http.StatusInternalServerError).Handle(ctx, w)
		return
	}
}

// LINE Notify の連携を解除する
func LinenotifyOff(w http.ResponseWriter, req *http.Request) {
	// Using feature that provided in official.
	http.Redirect(w, req, "https://notify-bot.line.me/my/", http.StatusFound)
}

func LinenotifyCallback(w http.ResponseWriter, req *http.Request) {
	ctx := getContext(req)

	params, err := linenotify.ParseCallbackParameters(req)
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

	reqToken := linenotify.NewRequestToken(
		params.Code,
		buildURL(req.URL, "/linenotify/callback"),
		os.Getenv("LINENOTIFY_CLIENT_ID"),
		os.Getenv("LINENOTIFY_CLIENT_SECRET"),
	)
	reqToken.Client = urlfetch.Client(ctx)

	token, err := reqToken.Get()
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
}
