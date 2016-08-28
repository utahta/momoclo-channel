package linebot

import (
	"net"
	"time"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/appengine/lib/log"
	"github.com/utahta/momoclo-channel/appengine/model"
	"github.com/utahta/momoclo-channel/crawler"
	pb "github.com/utahta/momoclo-channel/linebot/protos"
	"golang.org/x/net/context"
	"google.golang.org/appengine/socket"
	"google.golang.org/grpc"
)

type Client struct {
	context       context.Context
	conn          *grpc.ClientConn
	LineBotClient pb.LineBotClient
	Log           log.Logger
}

func Dial(ctx context.Context, address string) (*Client, error) {
	dialOption := grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
		return socket.DialTimeout(ctx, "tcp", addr, timeout)
	})
	conn, err := grpc.Dial(address, grpc.WithInsecure(), dialOption)
	if err != nil {
		return nil, errors.Wrapf(err, "did not connect. address:%s", address)
	}
	return &Client{context: ctx, conn: conn, LineBotClient: pb.NewLineBotClient(conn), Log: log.NewGaeLogger(ctx)}, nil
}

func (c *Client) Close() {
	if c.conn == nil {
		return
	}
	if err := c.conn.Close(); err != nil {
		c.Log.Errorf("Failed to close connection. error:%v", err)
	}
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
	// make gRPC request params.
	reqItem := &pb.NotifyChannelRequest_Item{Title: item.Title, Url: item.Url}
	for _, image := range item.Images {
		reqItem.Images = append(reqItem.Images, &pb.NotifyChannelRequest_Item_Image{Url: image.Url})
	}
	for _, video := range item.Videos {
		reqItem.Videos = append(reqItem.Videos, &pb.NotifyChannelRequest_Item_Video{Url: video.Url})
	}
	req := &pb.NotifyChannelRequest{Title: title, Item: reqItem}

	c.notifyAll(func(ctx context.Context, to []string) error {
		req.To = to
		if _, err := c.LineBotClient.NotifyChannel(ctx, req); err != nil {
			return errors.Wrapf(err, "Failed to notify channel. title:%s", item.Title)
		}
		c.Log.Infof("Notify channel. title:%s", item.Title)
		return nil
	})
	return nil
}
