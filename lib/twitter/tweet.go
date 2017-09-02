package twitter

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/utahta/go-twitter/types"
	"github.com/utahta/momoclo-channel/lib/config"
	"github.com/utahta/momoclo-channel/lib/log"
	"github.com/utahta/momoclo-crawler"
)

type ChannelParam struct {
	Title string
	Item  *crawler.ChannelItem
}

// Tweet text message
func TweetMessage(ctx context.Context, text string) error {
	if config.C.Twitter.Disabled {
		return nil
	}

	c, err := newClient(ctx)
	if err != nil {
		return err
	}

	if _, err := c.Tweet(text, nil); err != nil {
		return err
	}
	return nil
}

// Tweet channel
func TweetChannel(ctx context.Context, param *ChannelParam) error {
	if config.C.Twitter.Disabled {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, 540*time.Second)
	defer cancel()

	c, err := newClient(ctx)
	if err != nil {
		return err
	}

	const maxUploadMediaLen = 4
	var images [][]string
	var tmp []string
	for _, image := range param.Item.Images {
		tmp = append(tmp, image.Url)
		if len(tmp) == maxUploadMediaLen {
			images = append(images, tmp)
			tmp = nil
		}
	}
	if len(tmp) > 0 {
		images = append(images, tmp)
	}
	videos := param.Item.Videos
	text := truncateText(param.Title, param.Item)

	var tweets *types.Tweets
	if len(images) > 0 {
		tweets, err = c.TweetImageURLs(text, images[0], nil)
		images = images[1:]
	} else if len(videos) > 0 {
		tweets, err = c.TweetVideoURL(text, videos[0].Url, "video/mp4", nil)
		videos = videos[1:]
	} else {
		tweets, err = c.Tweet(text, nil)
	}

	if err != nil {
		log.Errorf(ctx, "Failed to post tweet. text:%s err:%v", text, err)
		return err
	}
	log.Infof(ctx, "Post tweet. text:%s images:%v videos:%v", text, len(param.Item.Images), len(param.Item.Videos))

	if len(images) > 0 {
		for _, urlsStr := range images {
			v := url.Values{}
			v.Set("in_reply_to_status_id", tweets.IDStr)
			tweets, err = c.TweetImageURLs("", urlsStr, v)
			if err != nil {
				log.Errorf(ctx, "Failed to post images. urls:%v err:%v", urlsStr, err)
			}
		}
	}

	if len(videos) > 0 {
		for _, video := range videos {
			v := url.Values{}
			v.Set("in_reply_to_status_id", tweets.IDStr)
			tweets, err = c.TweetVideoURL("", video.Url, "video/mp4", v)
			if err != nil {
				log.Errorf(ctx, "Failed to post video. url:%v err:%v", video.Url, err)
			}
		}
	}
	return nil
}

func truncateText(channelTitle string, item *crawler.ChannelItem) string {
	const maxTweetTextLen = 77 // ハッシュタグや URL や画像を除いて投稿可能な文字数

	title := []rune(fmt.Sprintf("%s %s", channelTitle, item.Title))
	if len(title) >= maxTweetTextLen {
		title = append(title[0:maxTweetTextLen-3], []rune("...")...)
	}
	return fmt.Sprintf("%s %s #momoclo #ももクロ", string(title), item.Url)
}
