package backend

import (
	"html/template"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/utahta/go-linenotify/auth"
	"github.com/utahta/momoclo-channel/adapter/handler"
	"github.com/utahta/momoclo-channel/container"
	"github.com/utahta/momoclo-channel/lib/config"
	"github.com/utahta/momoclo-channel/usecase"
)

// LineNotifyOn redirect to LINE Notify connection page
func LineNotifyOn(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	c, err := auth.New(config.C.LineNotify.ClientID, config.LineNotifyCallbackURL())
	if err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "state", Value: c.State, Expires: time.Now().Add(300 * time.Second), Secure: true})

	container.Logger(ctx).AE().Info(ctx, "Redirect to LINE Notify connection page")

	err = c.Redirect(w, req)
	if err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}
}

// LineNotifyOff redirect to LINE Notify revoking page
func LineNotifyOff(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	container.Logger(ctx).AE().Info(ctx, "Redirect to LINE Notify revoking page")

	// official url
	http.Redirect(w, req, "https://notify-bot.line.me/my/", http.StatusFound)
}

// LineNotifyCallback stores LINE Notify token
func LineNotifyCallback(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	params, err := auth.ParseRequest(req)
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

	if err := container.Usecase(ctx).AddLineNotification().Do(usecase.AddLineNotificationParams{Code: params.Code}); err != nil {
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
}
