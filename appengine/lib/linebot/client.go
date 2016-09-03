package linebot

import (
	"fmt"
	"os"
	"strconv"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/appengine/lib/log"
	"github.com/utahta/momoclo-channel/appengine/model"
	"github.com/utahta/momoclo-channel/crawler"
	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
)

type Client struct {
	context       context.Context
	LineBotClient *linebot.Client
	Log           log.Logger
}

func NewClient(ctx context.Context) (*Client, error) {
	var (
		channelID     int64
		channelSecret = os.Getenv("LINEBOT_CHANNEL_SECRET")
		channelMID    = os.Getenv("LINEBOT_CHANNEL_MID")
	)
	channelID, err := strconv.ParseInt(os.Getenv("LINEBOT_CHANNEL_ID"), 10, 64)
	if err != nil {
		return nil, err
	}
	bot, err := linebot.NewClient(channelID, channelSecret, channelMID, linebot.WithHTTPClient(urlfetch.Client(ctx)))
	if err != nil {
		return nil, err
	}
	return &Client{context: ctx, LineBotClient: bot, Log: log.NewGaeLogger(ctx)}, nil
}

func (c *Client) notifyAll(fn func(context.Context, []string) error) error {
	var (
		err error
		to  []string
		q   = model.NewLineUserQuery(c.context)
	)
	for {
		to, err = q.GetIds()
		if err != nil {
			return errors.Wrapf(err, "Failed to get user ids.")
		}
		count := len(to)

		if count > 0 {
			if err := fn(c.context, to); err != nil {
				return err
			}
		}
		if count < q.Limit {
			break
		}
	}
	return nil
}

func (c *Client) NotifyChannel(title string, item *crawler.ChannelItem) error {
	req := c.LineBotClient.NewMultipleMessage()
	req.AddText(fmt.Sprintf("%s\n%s\n%s", title, item.Title, item.Url))
	for _, image := range item.Images {
		req.AddImage(image.Url, image.Url)
	}

	return c.notifyAll(func(ctx context.Context, to []string) error {
		if _, err := req.Send(to); err != nil {
			return errors.Wrapf(err, "Failed to notify channel. title:%s", item.Title)
		}
		c.Log.Infof("Notify channel. title:%s count:%d", item.Title, len(to))
		return nil
	})
}

func (c *Client) NotifyMessage(text string) error {
	return c.notifyAll(func(ctx context.Context, to []string) error {
		if _, err := c.LineBotClient.SendText(to, text); err != nil {
			return errors.Wrap(err, "Failed to notify message.")
		}
		c.Log.Infof("Notify message. count:%d", len(to))
		return nil
	})
}

func (c *Client) NotifyMessageTo(to []string, text string) error {
	if _, err := c.LineBotClient.SendText(to, text); err != nil {
		return errors.Wrap(err, "Failed to notify message to.")
	}
	c.Log.Infof("Notify message to. count:%d", len(to))
	return nil
}
