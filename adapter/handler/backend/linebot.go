package backend

import (
	"html/template"
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/utahta/momoclo-channel/adapter/handler"
	mbot "github.com/utahta/momoclo-channel/lib/linebot"
)

func LineBotCallback(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	events, err := mbot.ParseRequest(ctx, req)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			handler.Fail(ctx, w, err, http.StatusBadRequest)
			return
		}
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}

	if err := mbot.HandleEvents(ctx, events); err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}
}

func LineBotHelp(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	tpl := template.Must(template.ParseFiles("view/linebot/help.html"))
	if err := tpl.Execute(w, nil); err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}
}

func LineBotAbout(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	tpl := template.Must(template.ParseFiles("view/linebot/about.html"))
	if err := tpl.Execute(w, nil); err != nil {
		handler.Fail(ctx, w, err, http.StatusInternalServerError)
		return
	}
}
