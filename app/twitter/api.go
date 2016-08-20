package twitter

import (
	"io/ioutil"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"github.com/ChimeraCoder/anaconda"
	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
	"github.com/utahta/momoclo-channel/crawler"
	"github.com/utahta/momoclo-channel/app/log"
	"github.com/pkg/errors"
)

const (
	maxUploadImages = 4
)

type TwitterApi struct {
	api *anaconda.TwitterApi
	Context context.Context
}

type mediaImage struct {
	Ids [maxUploadImages]string
}

func (t *TwitterApi) Auth(consumerKey, consumerSecret, accessToken, accessTokenSecret string) {
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	t.api = anaconda.NewTwitterApi(accessToken, accessTokenSecret)
	if t.Context != nil {
		t.api.HttpClient.Transport = &urlfetch.Transport{ Context: t.Context }
	}
}

func (t *TwitterApi) Tweet(ch *crawler.Channel) {
	for _, item := range ch.Items {
		images := t.uploadImages(item)
		videos := t.uploadVideos(item)

		text := t.truncateText(ch, item)
		v := url.Values{}
		if len(images) > 0 {
			v.Add("media_ids", strings.Join(images[0].Ids[:], ","))
			images = images[1:]
		} else if len(videos) > 0 {
			v.Add("media_ids", videos[0].MediaIDString)
			videos = videos[1:]
		}
		tweet, err := t.api.PostTweet(text, v)
		if err != nil {
			log.Errorf(t.Context, "Failed to post tweet. url:%s error:%s", item.Url, err)
			continue
		}
		log.Infof(t.Context, "Post tweet. text:%s", text)

		for _, image := range images {
			v := url.Values{}
			v.Add("in_reply_to_status_id", tweet.IdStr)
			v.Add("media_ids", strings.Join(image.Ids[:], ","))

			tweet, err = t.api.PostTweet("", v)
			if err != nil {
				log.Errorf(t.Context, "Failed to post tweet images. error:%v", err)
				continue
			}
		}

		for _, video := range videos {
			v := url.Values{}
			v.Add("in_reply_to_status_id", tweet.IdStr)
			v.Add("media_ids", video.MediaIDString)

			tweet, err = t.api.PostTweet("", v)
			if err != nil {
				log.Errorf(t.Context, "Failed to post tweet videos. error:%v", err)
				continue
			}
		}
	}
}

func (t *TwitterApi) truncateText(ch *crawler.Channel, item *crawler.ChannelItem) string {
	const maxTweetTextLen = 101 // ハッシュタグや url を除いて投稿可能な文字数

	title := []rune(fmt.Sprintf("%s %s", ch.Title, item.Title))
	if len(title) > maxTweetTextLen {
		title = append(title[0:maxTweetTextLen-3], []rune("...")...)
	}
	return fmt.Sprintf("%s %s #momoclo #ももクロ", string(title), item.Url)
}

func (t *TwitterApi) uploadImages(item *crawler.ChannelItem) ([]*mediaImage) {
	ids := []string{}
	for _, image := range item.Images {
		resource, err := t.downloadImage(image.Url)
		if err != nil {
			log.Errorf(t.Context, "url:%s error:%s", image.Url, err)
			continue
		}
		media, err := t.api.UploadMedia(resource)
		if err != nil {
			log.Errorf(t.Context, "Failed to upload media. url:%s error:%s", image.Url, err)
			continue
		}
		ids = append(ids, media.MediaIDString)
	}

	mis := []*mediaImage{}
	mi := &mediaImage{}
	num := 0
	for _, v := range ids {
		mi.Ids[num] = v
		num++
		// twitter max upload images
		if num == maxUploadImages {
			num = 0
			mis = append(mis, mi)
			mi = &mediaImage{}
		}
	}
	if mi.Ids[0] != "" {
		mis = append(mis, mi)
	}
	return mis
}

func (t *TwitterApi) uploadVideos(item *crawler.ChannelItem) ([]*anaconda.VideoMedia) {
	videos := []*anaconda.VideoMedia{}
	for _, video := range item.Videos {
		resp, err := t.api.HttpClient.Get(video.Url)
		if err != nil {
			log.Errorf(t.Context, "failed to get mp4 url:%s err:%v\n", video.Url, err)
			continue
		}

		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Errorf(t.Context, "failed to readall mp4 url:%s err:%v\n", video.Url, err)
			continue
		}

		media, err := t.api.UploadVideoInit(len(bytes), "video/mp4")
		if err != nil {
			log.Errorf(t.Context, "failed to upload video init. url:%s err:%v\n", video.Url, err)
			continue
		}
		if err = t.api.UploadVideoAppend(media.MediaIDString, 0, base64.StdEncoding.EncodeToString(bytes)); err != nil {
			log.Errorf(t.Context, "failed to upload video append. url:%s err:%v\n", video.Url, err)
			continue
		}
		v, err := t.api.UploadVideoFinalize(media.MediaIDString)
		if err != nil {
			log.Errorf(t.Context, "failed to upload video finalize. url:%s err:%v\n", video.Url, err)
			continue
		}
		videos = append(videos, &v)
	}
	return videos
}


func (t *TwitterApi) downloadImage(url string) (string, error) {
	response, err := t.api.HttpClient.Get(url)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to download image. url:%s", url)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to read response. url:%s", url)
	}
	return base64.StdEncoding.EncodeToString(body), nil
}
