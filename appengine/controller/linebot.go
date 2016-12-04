package controller

import (
	"html/template"
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
	mbot "github.com/utahta/momoclo-channel/appengine/lib/linebot"
	"golang.org/x/net/context"
)

func LineBotCallback(w http.ResponseWriter, req *http.Request) {
	ctx := getContext(req)

	events, err := mbot.ParseRequest(ctx, req)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			newError(err, http.StatusBadRequest).Handle(ctx, w)
			return
		}
		newError(err, http.StatusInternalServerError).Handle(ctx, w)
		return
	}

	ctx = context.WithValue(ctx, "baseURL", buildURL(req.URL, ""))
	if err := mbot.HandleEvents(ctx, events); err != nil {
		newError(err, http.StatusInternalServerError).Handle(ctx, w)
		return
	}
}

func LineBotHelp(w http.ResponseWriter, req *http.Request) {
	ctx := getContext(req)

	tpl := template.Must(template.ParseFiles("view/linebot/help.html"))
	if err := tpl.Execute(w, nil); err != nil {
		newError(err, http.StatusInternalServerError).Handle(ctx, w)
		return
	}
}

func LineBotAbout(w http.ResponseWriter, req *http.Request) {
	ctx := getContext(req)

	tpl := template.Must(template.ParseFiles("view/linebot/about.html"))
	if err := tpl.Execute(w, nil); err != nil {
		newError(err, http.StatusInternalServerError).Handle(ctx, w)
		return
	}
}
