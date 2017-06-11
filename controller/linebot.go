package controller

import (
	"html/template"
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/utahta/momoclo-channel/app"
	mbot "github.com/utahta/momoclo-channel/lib/linebot"
	"golang.org/x/net/context"
)

func LineBotCallback(w http.ResponseWriter, req *http.Request) {
	ctx := app.GetContext(req)

	events, err := mbot.ParseRequest(ctx, req)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			ctx.Error(err, http.StatusBadRequest)
			return
		}
		ctx.Fail(err)
		return
	}

	withCtx := context.WithValue(ctx, "baseURL", buildURL(req.URL, ""))
	if err := mbot.HandleEvents(withCtx, events); err != nil {
		ctx.Fail(err)
		return
	}
}

func LineBotHelp(w http.ResponseWriter, req *http.Request) {
	ctx := app.GetContext(req)

	tpl := template.Must(template.ParseFiles("view/linebot/help.html"))
	if err := tpl.Execute(w, nil); err != nil {
		ctx.Fail(err)
		return
	}
}

func LineBotAbout(w http.ResponseWriter, req *http.Request) {
	ctx := app.GetContext(req)

	tpl := template.Must(template.ParseFiles("view/linebot/about.html"))
	if err := tpl.Execute(w, nil); err != nil {
		ctx.Fail(err)
		return
	}
}
