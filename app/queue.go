package app

import (
	"encoding/json"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/crawler"
	"github.com/utahta/momoclo-channel/line/client"
	pb "github.com/utahta/momoclo-channel/line/protos"
	"github.com/utahta/momoclo-channel/log"
	"github.com/utahta/momoclo-channel/model"
	"github.com/utahta/momoclo-channel/twitter"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/socket"
	"google.golang.org/appengine/urlfetch"
	"google.golang.org/grpc"
)

// Queue for crawler.Channel
type QueueHandler struct {
	log log.Logger
}

func (h *QueueHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	h.log = NewGaeLogger(ctx)
	var err *Error

	ch, err := h.parseParams(r)
	if err != nil {
		err.Handle(ctx, w)
		return
	}

	switch r.URL.Path {
	case "/queue/tweet":
		err = h.tweet(ctx, ch)
	case "/queue/line":
		err = h.line(ctx, ch)
	default:
		http.NotFound(w, r)
	}
	err.Handle(ctx, w)
}

func (h *QueueHandler) parseParams(r *http.Request) (*crawler.Channel, *Error) {
	var ch crawler.Channel
	if err := json.Unmarshal([]byte(r.FormValue("channel")), &ch); err != nil {
		return nil, newError(errors.Wrapf(err, "Failed to unmarshal."), http.StatusInternalServerError)
	}
	return &ch, nil
}

func (h *QueueHandler) tweet(ctx context.Context, ch *crawler.Channel) *Error {
	ctx, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()
	client := urlfetch.Client(ctx)

	var wg sync.WaitGroup
	wg.Add(len(ch.Items))
	for _, item := range ch.Items {
		go func(ctx context.Context, item *crawler.ChannelItem) {
			defer wg.Done()

			if err := model.NewTweetItem(item).Put(ctx); err != nil {
				return
			}

			tw := twitter.NewChannelClient(
				os.Getenv("TWITTER_CONSUMER_KEY"),
				os.Getenv("TWITTER_CONSUMER_SECRET"),
				os.Getenv("TWITTER_ACCESS_TOKEN"),
				os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"),
			)
			tw.Log = NewGaeLogger(ctx)
			tw.Api.HttpClient = client

			tw.TweetItem(ch.Title, item)
		}(ctx, item)
	}
	wg.Wait()

	return nil
}

func (h *QueueHandler) line(ctx context.Context, ch *crawler.Channel) *Error {
	ctx, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	dialOption := grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
		return socket.DialTimeout(ctx, "tcp", os.Getenv("LINE_SERVER_ADDRESS"), timeout)
	})
	client, err := client.Dial(os.Getenv("LINE_SERVER_ADDRESS"), grpc.WithInsecure(), dialOption)
	if err != nil {
		return newError(err, http.StatusInternalServerError)
	}
	defer client.Close()
	client.Log = h.log

	var wg sync.WaitGroup
	wg.Add(len(ch.Items))
	for _, item := range ch.Items {
		go func(ctx context.Context, item *crawler.ChannelItem) {
			defer wg.Done()

			if err := model.NewLineItem(item).Put(ctx); err != nil {
				return
			}

			// make gRPC request params.
			reqItem := &pb.NotifyChannelRequest_Item{Title: item.Title, Url: item.Url}
			for _, image := range item.Images {
				reqItem.Images = append(reqItem.Images, &pb.NotifyChannelRequest_Item_Image{Url: image.Url})
			}
			for _, video := range item.Videos {
				reqItem.Videos = append(reqItem.Videos, &pb.NotifyChannelRequest_Item_Video{Url: video.Url})
			}
			req := &pb.NotifyChannelRequest{Title: ch.Title, Item: reqItem}

			// notify channel item
			var (
				q      = model.NewUserQuery(ctx)
				cursor = datastore.Cursor{}
				err    error
			)
			for {
				req.To, cursor, err = q.GetIds(cursor)
				if err != nil {
					h.log.Errorf("Failed to get user ids. error:%v", err)
					return
				}
				count := len(req.To)

				if count > 0 {
					if _, err := client.NotifyChannel(ctx, req); err != nil {
						h.log.Errorf("Failed to notify channel. error:%v", err)
						return
					}
					h.log.Infof("Notify channel. title:%s", req.Item.Title)
				}
				if count < q.Limit {
					break
				}
			}
		}(ctx, item)
	}
	wg.Wait()

	return nil
}
