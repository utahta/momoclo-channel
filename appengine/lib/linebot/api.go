package linebot

import (
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
	"golang.org/x/net/context"
)

func ParseRequest(ctx context.Context, req *http.Request) ([]*linebot.Event, error) {
	cli, err := New(ctx)
	if err != nil {
		return nil, err
	}
	bot := cli.LineBotClient

	events, err := bot.ParseRequest(req)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func ReplyText(ctx context.Context, replyToken, text string) error {
	client, err := New(ctx)
	if err != nil {
		return err
	}

	if err := client.ReplyText(replyToken, text); err != nil {
		return err
	}
	return nil
}

func ReplyImage(ctx context.Context, replyToken, originalContentURL, previewImageURL string) error {
	client, err := New(ctx)
	if err != nil {
		return err
	}

	if err := client.ReplyImage(replyToken, originalContentURL, previewImageURL); err != nil {
		return err
	}
	return nil
}
