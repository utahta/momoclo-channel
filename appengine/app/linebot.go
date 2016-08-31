package app

import (
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
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

	bot, err := h.newBotClient()
	if err != nil {
		return newError(err, http.StatusInternalServerError)
	}

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
			err := h.appendUser(ctx, content)
			if err != nil {
				h.log.Error(err)
				continue
			}
		} else if content.IsOperation && content.OpType == linebot.OpTypeBlocked {
			err := h.deleteUser(ctx, content)
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

			if err := h.parseText(ctx, text); err != nil {
				h.log.Error(err)
				continue
			}
		}
	}
	return nil
}

func (h *LinebotHandler) appendUser(ctx context.Context, content *linebot.ReceivedContent) error {
	user := model.NewLineUser(content.From)
	user.Enabled = true
	if err := user.Put(ctx); err != nil {
		return err
	}
	mbot.NotifyMessage(ctx, "通知ノフ設定オンにしました（・Θ・）")
	return nil
}

func (h *LinebotHandler) deleteUser(ctx context.Context, content *linebot.ReceivedContent) error {
	user := model.NewLineUser(content.From)
	user.Enabled = false
	if err := user.Put(ctx); err != nil {
		return err
	}
	mbot.NotifyMessage(ctx, "通知ノフ設定オフにしました（・Θ・）")
	return nil
}

func (h *LinebotHandler) parseText(ctx context.Context, text *linebot.ReceivedTextContent) error {
	var (
		matched bool
		err     error
	)
	matched, err = regexp.MatchString("^(おん|オン|on)$", text.Text)
	if err != nil {
		return err
	}
	if matched {
		return h.appendUser(ctx, text.ReceivedContent)
	}

	matched, err = regexp.MatchString("^(おふ|オフ|off)$", text.Text)
	if err != nil {
		return err
	}
	if matched {
		return h.deleteUser(ctx, text.ReceivedContent)
	}

	mbot.NotifyMessage(ctx, "?（・Θ・）?")
	return nil
}

func (h *LinebotHandler) newBotClient() (*linebot.Client, error) {
	var (
		channelID     int64
		channelSecret = os.Getenv("LINEBOT_CHANNEL_SECRET")
		channelMID    = os.Getenv("LINEBOT_CHANNEL_MID")
	)
	channelID, err := strconv.ParseInt(os.Getenv("LINEBOT_CHANNEL_ID"), 10, 64)
	if err != nil {
		return nil, err
	}
	return linebot.NewClient(channelID, channelSecret, channelMID)
}
