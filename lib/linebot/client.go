package linebot

import (
	"context"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/utahta/momoclo-channel/lib/config"
	"github.com/utahta/momoclo-channel/lib/log"
	"google.golang.org/appengine/urlfetch"
)

type Client struct {
	context       context.Context
	LineBotClient *linebot.Client
}

func New(ctx context.Context) (*Client, error) {
	var (
		channelSecret = config.C.LineBot.ChannelSecret
		channelToken  = config.C.LineBot.ChannelToken
	)
	bot, err := linebot.New(channelSecret, channelToken, linebot.WithHTTPClient(urlfetch.Client(ctx)))
	if err != nil {
		return nil, err
	}
	return &Client{context: ctx, LineBotClient: bot}, nil
}

func (c *Client) ReplyText(replyToken, text string) error {
	if _, err := c.LineBotClient.ReplyMessage(
		replyToken,
		linebot.NewTextMessage(text),
	).Do(); err != nil {
		return err
	}
	log.Infof(c.context, "Reply text. text:%s", text)
	return nil
}

func (c *Client) ReplyImage(replyToken, originalContentURL, previewImageURL string) error {
	if _, err := c.LineBotClient.ReplyMessage(
		replyToken,
		linebot.NewImageMessage(originalContentURL, previewImageURL),
	).Do(); err != nil {
		return err
	}
	log.Info(c.context, "Reply image.")
	return nil
}
