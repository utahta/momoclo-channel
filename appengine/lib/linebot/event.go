package linebot

import (
	"fmt"
	"regexp"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/appengine/lib/googleapi/customsearch"
	"github.com/utahta/momoclo-channel/appengine/lib/log"
	"golang.org/x/net/context"
)

var (
	ErrorHandleOnOffNotMatch = errors.New("handle on off not match.")
	ErrorHandleImageNotMatch = errors.New("handle image not match.")
)

// handle line message events.
func HandleEvents(ctx context.Context, events []*linebot.Event) error {
	for _, event := range events {
		switch event.Type {
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if err := textMessageEvent(ctx, message, event); err != nil {
					log.GaeLog(ctx).Error(err)
					continue
				}
			}
		case linebot.EventTypeFollow:
			err := followEvent(ctx, event)
			if err != nil {
				return err
			}
		case linebot.EventTypeUnfollow:
			err := unfollowEvent(ctx, event)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func onMessage(ctx context.Context) string {
	onURL := fmt.Sprintf("%s%s", ctx.Value("baseURL").(string), "/linenotify/on")
	return fmt.Sprintf("通知機能を有効にする場合は、下記URLから設定を行ってください（・Θ・）\n%s", onURL)
}

func helpMessage(ctx context.Context) string {
	helpURL := fmt.Sprintf("%s%s", ctx.Value("baseURL").(string), "/linebot/help")
	return fmt.Sprintf("ヘルプ（・Θ・）\n%s", helpURL)
}

func followEvent(ctx context.Context, event *linebot.Event) error {
	log.GaeLog(ctx).Infof("follow user. event:%v", event)

	message := "友だち追加ありがとうございます。\nこちらは、ももクロちゃんのブログやAE NEWS等を通知する機能と連携したり、画像を返したりするBOTです。"
	return ReplyText(ctx, event.ReplyToken, fmt.Sprintf("%s\n\n%s\n\n%s", message, onMessage(ctx), helpMessage(ctx)))
}

func unfollowEvent(ctx context.Context, event *linebot.Event) error {
	log.GaeLog(ctx).Infof("unfollow user. event:%v", event)
	return nil
}

func textMessageEvent(ctx context.Context, message *linebot.TextMessage, event *linebot.Event) error {
	log.GaeLog(ctx).Infof("handle text content. message:%s", message.Text)

	if err := handleOnOff(ctx, message, event); err != ErrorHandleOnOffNotMatch {
		return err
	}

	if err := handleMemberImage(ctx, message, event); err != ErrorHandleImageNotMatch {
		return err
	}

	return ReplyText(ctx, event.ReplyToken, helpMessage(ctx))
}

func handleOnOff(ctx context.Context, message *linebot.TextMessage, event *linebot.Event) error {
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
		return ReplyText(ctx, event.ReplyToken, onMessage(ctx))
	}

	matched, err = regexp.MatchString("^(おふ|オフ|off)$", text)
	if err != nil {
		return err
	}
	if matched {
		return ReplyText(ctx, event.ReplyToken, "通知機能を無効にする場合は、下記URLから解除を行ってください（・Θ・）\nhttps://notify-bot.line.me/my/")
	}

	return ErrorHandleOnOffNotMatch
}

func handleMemberImage(ctx context.Context, message *linebot.TextMessage, event *linebot.Event) error {
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
		return ErrorHandleImageNotMatch
	}

	res, err := customsearch.SearchImage(ctx, word)
	if err != nil {
		log.GaeLog(ctx).Error(err)
		return ReplyText(ctx, event.ReplyToken, "画像がみつかりませんでした（・Θ・）")
	}
	return ReplyImage(ctx, event.ReplyToken, res.Url, res.ThumbnailUrl)
}
