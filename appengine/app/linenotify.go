package app

import (
	"html/template"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/appengine/lib/log"
	"github.com/utahta/momoclo-channel/linenotify"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

type LinenotifyHandler struct {
	log log.Logger
}

func (h *LinenotifyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(appengine.NewContext(r), 50*time.Second)
	defer cancel()

	h.log = log.NewGaeLogger(ctx)
	var err *Error

	switch r.URL.Path {
	case "/linenotify/on":
		err = h.handleOn(ctx, w, r)
	case "/linenotify/off":
		err = h.handleOff(ctx, w, r)
	case "/linenotify/callback":
		err = h.handleCallback(ctx, w, r)
	default:
		http.NotFound(w, r)
	}
	err.Handle(ctx, w)
}

func callbackURI(req *http.Request) (string, error) {
	url, err := url.Parse(req.URL.String())
	if err != nil {
		return "", err
	}
	url.Path = "/linenotify/callback"
	return url.String(), nil
}

// LINE Notify と連携する
func (h *LinenotifyHandler) handleOn(ctx context.Context, w http.ResponseWriter, req *http.Request) *Error {
	uri, err := callbackURI(req)
	if err != nil {
		return newError(err, http.StatusInternalServerError)
	}

	reqAuth, err := linenotify.NewRequestAuthorization(os.Getenv("LINENOTIFY_CLIENT_ID"), uri)
	if err != nil {
		return newError(err, http.StatusInternalServerError)
	}
	http.SetCookie(w, &http.Cookie{Name: "state", Value: reqAuth.State, Expires: time.Now().Add(60 * time.Second), Secure: true})

	err = reqAuth.Redirect(w, req)
	if err != nil {
		return newError(err, http.StatusInternalServerError)
	}

	return nil
}

// LINE Notify の連携を解除する
func (h *LinenotifyHandler) handleOff(ctx context.Context, w http.ResponseWriter, req *http.Request) *Error {
	// Using feature that provided in official.
	http.Redirect(w, req, "https://notify-bot.line.me/my/", http.StatusFound)
	return nil
}

// LINE Notify 連携ページからの認証パラメータを受付
func (h *LinenotifyHandler) handleCallback(ctx context.Context, w http.ResponseWriter, req *http.Request) *Error {
	params, err := linenotify.ParseCallbackParameters(req)
	if err != nil {
		return newError(err, http.StatusInternalServerError)
	}

	state, err := req.Cookie("state")
	if err != nil {
		return newError(err, http.StatusInternalServerError)
	}

	if params.State != state.Value {
		return newError(errors.New("Invalid csrf token."), http.StatusBadRequest)
	}

	uri, err := callbackURI(req)
	if err != nil {
		return newError(err, http.StatusInternalServerError)
	}

	reqToken := linenotify.NewRequestToken(params.Code, uri, os.Getenv("LINENOTIFY_CLIENT_ID"), os.Getenv("LINENOTIFY_CLIENT_SECRET"))
	reqToken.Client = urlfetch.Client(ctx)

	_, err = reqToken.Get()
	if err != nil {
		return newError(err, http.StatusInternalServerError)
	}

	t, err := template.New("callback").Parse("<html><body><h1>通知ノフ設定オンにしました（・Θ・）</h1></body></html>")
	if err != nil {
		return newError(err, http.StatusInternalServerError)
	}
	err = t.Execute(w, nil)
	if err != nil {
		return newError(err, http.StatusInternalServerError)
	}
	return nil
}
