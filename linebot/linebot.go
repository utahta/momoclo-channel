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
		ReplyText(string, string) error
		ReplyImage(string, string, string) error
	}

	// client represents line bot client
	client struct {
		*linebot.Client
	}
)

// New returns Client
func New(ctx context.Context) Client {
	c, _ := linebot.New(
		config.C().LineBot.ChannelSecret,
		config.C().LineBot.ChannelToken,
		linebot.WithHTTPClient(urlfetch.Client(ctx)),
	)
	return &client{c}
}

// ReplyText reply text message to bot
func (c *client) ReplyText(replyToken, text string) error {
	textMessage := linebot.NewTextMessage(text)
	if _, err := c.ReplyMessage(replyToken, textMessage).Do(); err != nil {
		return err
	}
	return nil
}

// ReplyImage reply image message to bot
func (c *client) ReplyImage(replyToken, originalContentURL, previewImageURL string) error {
	imageMessage := linebot.NewImageMessage(originalContentURL, previewImageURL)
	if _, err := c.ReplyMessage(replyToken, imageMessage).Do(); err != nil {
		return err
	}
	return nil
}
