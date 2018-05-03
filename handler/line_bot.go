package handler

import (
	"html/template"
	"net/http"

	"github.com/utahta/momoclo-channel/container"
	"github.com/utahta/momoclo-channel/linebot"
	"github.com/utahta/momoclo-channel/types"
	"github.com/utahta/momoclo-channel/usecase"
)

// LineBotCallback handler
func LineBotCallback(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	events, err := linebot.ParseRequest(req)
	if err != nil {
		if err == types.ErrInvalidSignature {
			failResponse(ctx, w, err, http.StatusBadRequest)
			return
		}
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}

	params := usecase.HandleLineBotEventsParams{Events: events}
	if err := container.Usecase(ctx).HandleLineBotEvents().Do(params); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}
}

// LineBotHelp handler
func LineBotHelp(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	tpl := template.Must(template.ParseFiles("public/templates/linebot/help.html"))
	if err := tpl.Execute(w, nil); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}
}

// LineBotAbout handler
func LineBotAbout(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	tpl := template.Must(template.ParseFiles("public/templates/linebot/about.html"))
	if err := tpl.Execute(w, nil); err != nil {
		failResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}
}
