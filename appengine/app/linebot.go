package app

import (
	"net/http"
	"regexp"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/utahta/momoclo-channel/appengine/lib/googleapi/customsearch"
	mbot "github.com/utahta/momoclo-channel/appengine/lib/linebot"
	"github.com/utahta/momoclo-channel/appengine/lib/log"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
)

type LinebotHandler struct {
	log log.Logger
}

func (h *LinebotHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	h.log = log.NewGaeLogger(ctx)
	var err *Error

	switch r.URL.Path {
	case "/linebot/callback":
		err = h.callback(ctx, r)
	default:
		http.NotFound(w, r)
	}
	err.Handle(ctx, w)
}

func (h *LinebotHandler) callback(ctx context.Context, req *http.Request) *Error {
	ctx, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	cli, err := mbot.New(ctx)
	if err != nil {
		return newError(err, http.StatusInternalServerError)
	}
	bot := cli.LineBotClient

	events, err := bot.ParseRequest(req)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			return newError(err, http.StatusBadRequest)
		}
		return newError(err, http.StatusInternalServerError)
	}

	for _, event := range events {
		switch event.Type {
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if err := h.handleTextMessage(ctx, message, event); err != nil {
					h.log.Error(err)
					continue
				}
			}
		case linebot.EventTypeFollow:
			err := h.followUser(ctx, event)
			if err != nil {
				return newError(err, http.StatusInternalServerError)
			}
		case linebot.EventTypeUnfollow:
			err := h.unfollowUser(ctx, event)
			if err != nil {
				return newError(err, http.StatusInternalServerError)
			}
		}
	}
	return nil
}

func (h *LinebotHandler) followUser(ctx context.Context, event *linebot.Event) error {
	h.log.Infof("follow user. event:%v", event)
	return mbot.ReplyText(ctx, event.ReplyToken, "通知ノフ設定オンにしました（・Θ・）")
}

func (h *LinebotHandler) unfollowUser(ctx context.Context, event *linebot.Event) error {
	h.log.Infof("unfollow user. event:%v", event)
	return mbot.ReplyText(ctx, event.ReplyToken, "通知ノフ設定オフにしました（・Θ・）")
}

func (h *LinebotHandler) handleTextMessage(ctx context.Context, message *linebot.TextMessage, event *linebot.Event) error {
	h.log.Infof("handle text content. event:%v", event)

	if err := h.handleOnOff(ctx, message, event); err != nil {
		return err
	}

	if err := h.handleMemberImage(ctx, message, event); err != nil {
		return err
	}

	return mbot.ReplyText(ctx, event.ReplyToken, "?（・Θ・）?\nヘルプ\nhttps://utahta.github.io/momoclo-channel/linebot/")
}

func (h *LinebotHandler) handleOnOff(ctx context.Context, message *linebot.TextMessage, event *linebot.Event) error {
	var (
		matched bool
		err     error
	)
	text := message.Text
	matched, err = regexp.MatchString("^(おん|オン|on)$", text)
	if err != nil {
		return err
	}
	if matched {
		return h.followUser(ctx, event)
	}

	matched, err = regexp.MatchString("^(おふ|オフ|off)$", text)
	if err != nil {
		return err
	}
	if matched {
		return h.unfollowUser(ctx, event)
	}

	return nil
}

func (h *LinebotHandler) handleMemberImage(ctx context.Context, message *linebot.TextMessage, event *linebot.Event) error {
	var (
		matched bool
		err     error
	)
	text := message.Text
	word := ""
	matched, err = regexp.MatchString("玉井|[たタ][まマ][いイ]|[しシ][おオ][りリ][んン]?|詩織|玉さん|[たタ][まマ]さん", text)
	if err != nil {
		return err
	}
	if matched {
		word = "玉井詩織"
	}
	matched, err = regexp.MatchString("百田|[もモ][もモ][たタ]|[夏かカ][菜なナ][子こコ]", text)
	if err != nil {
		return err
	}
	if matched {
		word = "百田夏菜子"
	}
	matched, err = regexp.MatchString("有安|[あア][りリ][やヤ][すス]|[もモ][もモ][かカ]|杏果", text)
	if err != nil {
		return err
	}
	if matched {
		word = "有安杏果"
	}
	matched, err = regexp.MatchString("佐々木|[さサ][さサ][きキ]|[あア][やヤ][かカ]|彩夏|[あア]ー[りリ][んン]", text)
	if err != nil {
		return err
	}
	if matched {
		word = "佐々木彩夏"
	}
	matched, err = regexp.MatchString("高城|[たタ][かカ][ぎギ]|[れレ][にニ]", text)
	if err != nil {
		return err
	}
	if matched {
		word = "高城れに"
	}

	if word == "" {
		return nil
	}

	res, err := customsearch.SearchImage(ctx, word)
	if err != nil {
		return err
	}
	return mbot.ReplyImage(ctx, event.ReplyToken, res.Url, res.ThumbnailUrl)
}
