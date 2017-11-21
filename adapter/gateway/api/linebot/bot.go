package linebot

import (
	"context"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/utahta/momoclo-channel/domain/model"
	"github.com/utahta/momoclo-channel/lib/config"
	"google.golang.org/appengine/urlfetch"
)

type (
	bot struct {
		*linebot.Client
	}
)

// New returns LineBot
func New(ctx context.Context) model.LineBot {
	c, _ := linebot.New(
		config.C.LineBot.ChannelSecret,
		config.C.LineBot.ChannelToken,
		linebot.WithHTTPClient(urlfetch.Client(ctx)),
	)
	return &bot{c}
}

// ReplyText reply text message to bot
func (c *bot) ReplyText(replyToken, text string) error {
	textMessage := linebot.NewTextMessage(text)
	if _, err := c.ReplyMessage(replyToken, textMessage).Do(); err != nil {
		return err
	}
	return nil
}

// ReplyImage reply image message to bot
func (c *bot) ReplyImage(replyToken, originalContentURL, previewImageURL string) error {
	imageMessage := linebot.NewImageMessage(originalContentURL, previewImageURL)
	if _, err := c.ReplyMessage(replyToken, imageMessage).Do(); err != nil {
		return err
	}
	return nil
}
