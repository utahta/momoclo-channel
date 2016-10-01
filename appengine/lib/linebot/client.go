package linebot

import (
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/utahta/momoclo-channel/appengine/lib/log"
	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
)

type Client struct {
	context       context.Context
	LineBotClient *linebot.Client
	Log           log.Logger
}

func New(ctx context.Context) (*Client, error) {
	var (
		channelSecret = os.Getenv("LINEBOT_CHANNEL_SECRET")
		channelToken  = os.Getenv("LINEBOT_CHANNEL_TOKEN")
	)
	bot, err := linebot.New(channelSecret, channelToken, linebot.WithHTTPClient(urlfetch.Client(ctx)))
	if err != nil {
		return nil, err
	}
	return &Client{context: ctx, LineBotClient: bot, Log: log.NewGaeLogger(ctx)}, nil
}

func (c *Client) ReplyText(replyToken, text string) error {
	if _, err := c.LineBotClient.ReplyMessage(
		replyToken,
		linebot.NewTextMessage(text),
	).Do(); err != nil {
		return err
	}
	c.Log.Infof("Reply text. text:%s", text)
	return nil
}

func (c *Client) ReplyImage(replyToken, originalContentURL, previewImageURL string) error {
	if _, err := c.LineBotClient.ReplyMessage(
		replyToken,
		linebot.NewImageMessage(originalContentURL, previewImageURL),
	).Do(); err != nil {
		return err
	}
	c.Log.Infof("Reply image.")
	return nil
}
