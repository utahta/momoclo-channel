package linebot

import (
	"context"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/utahta/momoclo-channel/config"
	"google.golang.org/appengine/urlfetch"
)

type (
	// Client represents line bot client interface
	Client interface {
		ReplyText(context.Context, string, string) error
		ReplyImage(context.Context, string, string, string) error
	}

	// client represents line bot client
	client struct {
	}
)

// New returns Client
func New() Client {
	return &client{}
}

// ReplyText reply text message to bot
func (c *client) ReplyText(ctx context.Context, replyToken, text string) error {
	bot, err := c.fromContext(ctx)
	if err != nil {
		return err
	}

	textMessage := linebot.NewTextMessage(text)
	if _, err := bot.ReplyMessage(replyToken, textMessage).Do(); err != nil {
		return err
	}
	return nil
}

// ReplyImage reply image message to bot
func (c *client) ReplyImage(ctx context.Context, replyToken, originalContentURL, previewImageURL string) error {
	bot, err := c.fromContext(ctx)
	if err != nil {
		return err
	}

	imageMessage := linebot.NewImageMessage(originalContentURL, previewImageURL)
	if _, err := bot.ReplyMessage(replyToken, imageMessage).Do(); err != nil {
		return err
	}
	return nil
}

func (c *client) fromContext(ctx context.Context) (*linebot.Client, error) {
	return linebot.New(
		config.C().LineBot.ChannelSecret,
		config.C().LineBot.ChannelToken,
		linebot.WithHTTPClient(urlfetch.Client(ctx)),
	)
}
