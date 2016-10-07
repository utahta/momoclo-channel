package controller

import (
	"html/template"
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
	mbot "github.com/utahta/momoclo-channel/appengine/lib/linebot"
	"golang.org/x/net/context"
)

func LineBotCallback(ctx context.Context, w http.ResponseWriter, req *http.Request) *Error {
	events, err := mbot.ParseRequest(ctx, req)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			return newError(err, http.StatusBadRequest)
		}
		return newError(err, http.StatusInternalServerError)
	}

	ctx = context.WithValue(ctx, "baseURL", buildURL(req.URL, ""))
	if err := mbot.HandleEvents(ctx, events); err != nil {
		return newError(err, http.StatusInternalServerError)
	}
	return nil
}

func LineBotHelp(ctx context.Context, w http.ResponseWriter, req *http.Request) *Error {
	tpl := template.Must(template.ParseFiles("../view/linebot/help.html"))
	if err := tpl.Execute(w, nil); err != nil {
		return newError(err, http.StatusInternalServerError)
	}
	return nil
}
