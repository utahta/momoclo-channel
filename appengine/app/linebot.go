package app

import (
	"net/http"
	"regexp"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/utahta/momoclo-channel/appengine/lib/googleapi/customsearch"
	mbot "github.com/utahta/momoclo-channel/appengine/lib/linebot"
	"github.com/utahta/momoclo-channel/appengine/lib/log"
	"github.com/utahta/momoclo-channel/appengine/model"
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

	client, err := mbot.NewClient(ctx)
	if err != nil {
		return newError(err, http.StatusInternalServerError)
	}
	bot := client.LineBotClient

	received, err := bot.ParseRequest(req)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			return newError(err, http.StatusBadRequest)
		}
		return newError(err, http.StatusInternalServerError)
	}

	for _, result := range received.Results {
		content := result.Content()
		if content == nil {
			h.log.Error("Invalid content.")
			continue
		}

		if content.IsOperation && content.OpType == linebot.OpTypeAddedAsFriend {
			h.log.Infof("append user. from:%s", content.From)
			err := h.appendUser(ctx, content.From)
			if err != nil {
				h.log.Error(err)
				continue
			}
		} else if content.IsOperation && content.OpType == linebot.OpTypeBlocked {
			h.log.Infof("delete user. from:%s", content.From)
			err := h.deleteUser(ctx, content.From)
			if err != nil {
				h.log.Error(err)
				continue
			}
		} else if content.IsMessage && content.ContentType == linebot.ContentTypeText {
			text, err := content.TextContent()
			if err != nil {
				h.log.Error(err)
				continue
			}
			h.log.Infof("handle text content. from:%s text:%s ", text.From, text.Text)

			if ok, err := h.handleOnOff(ctx, text.From, text.Text); ok || err != nil {
				if err != nil {
					h.log.Error(err)
				}
				continue
			}

			if ok, err := h.handleMemberImage(ctx, text.From, text.Text); ok || err != nil {
				if err != nil {
					h.log.Error(err)
				}
				continue
			}

			mbot.NotifyMessageTo(ctx, []string{text.From}, "?（・Θ・）?")
		}
	}
	return nil
}

func (h *LinebotHandler) appendUser(ctx context.Context, from string) error {
	user := model.NewLineUser(from)
	user.Enabled = true
	if err := user.Put(ctx); err != nil {
		return err
	}
	mbot.NotifyMessageTo(ctx, []string{user.Id}, "通知ノフ設定オンにしました（・Θ・）")
	return nil
}

func (h *LinebotHandler) deleteUser(ctx context.Context, from string) error {
	user := model.NewLineUser(from)
	user.Get(ctx)
	user.Enabled = false
	if err := user.Put(ctx); err != nil {
		return err
	}
	mbot.NotifyMessageTo(ctx, []string{user.Id}, "通知ノフ設定オフにしました（・Θ・）")
	return nil
}

func (h *LinebotHandler) handleOnOff(ctx context.Context, from, text string) (bool, error) {
	var (
		matched bool
		err     error
	)
	matched, err = regexp.MatchString("^(おん|オン|on)$", text)
	if err != nil {
		return false, err
	}
	if matched {
		return true, h.appendUser(ctx, from)
	}

	matched, err = regexp.MatchString("^(おふ|オフ|off)$", text)
	if err != nil {
		return false, err
	}
	if matched {
		return true, h.deleteUser(ctx, from)
	}

	return false, nil
}

func (h *LinebotHandler) handleMemberImage(ctx context.Context, from, text string) (bool, error) {
	var (
		matched bool
		err     error
	)
	word := ""
	matched, err = regexp.MatchString("玉井|[たタ][まマ][いイ]|[しシ][おオ][りリ][んン]?|詩織|玉さん|[たタ][まマ]さん", text)
	if err != nil {
		return false, err
	}
	if matched {
		word = "玉井詩織"
	}
	matched, err = regexp.MatchString("百田|[もモ][もモ][たタ]|[夏かカ][菜なナ][子こコ]", text)
	if err != nil {
		return false, err
	}
	if matched {
		word = "百田夏菜子"
	}
	matched, err = regexp.MatchString("有安|[あア][りリ][やヤ][すス]|[もモ][もモ][かカ]|杏果", text)
	if err != nil {
		return false, err
	}
	if matched {
		word = "有安杏果"
	}
	matched, err = regexp.MatchString("佐々木|[さサ][さサ][きキ]|[あア][やヤ][かカ]|彩夏|[あア]ー[りリ][んン]", text)
	if err != nil {
		return false, err
	}
	if matched {
		word = "佐々木彩夏"
	}
	matched, err = regexp.MatchString("高城|[たタ][かカ][ぎギ]|[れレ][にニ]", text)
	if err != nil {
		return false, err
	}
	if matched {
		word = "高城れに"
	}

	if word == "" {
		return false, nil
	}

	res, err := customsearch.SearchImage(ctx, word)
	if err != nil {
		return false, err
	}

	mbot.NotifyImageTo(ctx, []string{from}, res.Url, res.ThumbnailUrl)
	return true, nil
}
