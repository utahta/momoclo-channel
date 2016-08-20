package twitter

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/ChimeraCoder/anaconda"
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/crawler"
	"github.com/utahta/momoclo-channel/log"
)

const (
	maxUploadImages = 4
)

type TwitterClient struct {
	Api *anaconda.TwitterApi
	Log log.Logger
}

type mediaImage struct {
	Ids [maxUploadImages]string
}

func NewTwitterClient(consumerKey, consumerSecret, accessToken, accessTokenSecret string) *TwitterClient {
	t := &TwitterClient{}
	t.auth(consumerKey, consumerSecret, accessToken, accessTokenSecret)
	t.Log = log.NewSilentLogger()
	return t
}

func (t *TwitterClient) auth(consumerKey, consumerSecret, accessToken, accessTokenSecret string) {
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	t.Api = anaconda.NewTwitterApi(accessToken, accessTokenSecret)
}

func (t *TwitterClient) Tweet(ch *crawler.Channel) {
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
		tweet, err := t.Api.PostTweet(text, v)
		if err != nil {
			t.Log.Errorf("Failed to post tweet. url:%s error:%s", item.Url, err)
			continue
		}
		t.Log.Infof("Post tweet. text:%s", text)

		for _, image := range images {
			v := url.Values{}
			v.Add("in_reply_to_status_id", tweet.IdStr)
			v.Add("media_ids", strings.Join(image.Ids[:], ","))

			tweet, err = t.Api.PostTweet("", v)
			if err != nil {
				t.Log.Errorf("Failed to post tweet images. error:%v", err)
				continue
			}
		}

		for _, video := range videos {
			v := url.Values{}
			v.Add("in_reply_to_status_id", tweet.IdStr)
			v.Add("media_ids", video.MediaIDString)

			tweet, err = t.Api.PostTweet("", v)
			if err != nil {
				t.Log.Errorf("Failed to post tweet videos. error:%v", err)
				continue
			}
		}
	}
}

func (t *TwitterClient) truncateText(ch *crawler.Channel, item *crawler.ChannelItem) string {
	const maxTweetTextLen = 101 // ハッシュタグや url を除いて投稿可能な文字数

	title := []rune(fmt.Sprintf("%s %s", ch.Title, item.Title))
	if len(title) > maxTweetTextLen {
		title = append(title[0:maxTweetTextLen-3], []rune("...")...)
	}
	return fmt.Sprintf("%s %s #momoclo #ももクロ", string(title), item.Url)
}

func (t *TwitterClient) uploadImages(item *crawler.ChannelItem) []*mediaImage {
	ids := []string{}
	for _, image := range item.Images {
		resource, err := t.downloadImage(image.Url)
		if err != nil {
			t.Log.Errorf("url:%s error:%s", image.Url, err)
			continue
		}
		media, err := t.Api.UploadMedia(resource)
		if err != nil {
			t.Log.Errorf("Failed to upload media. url:%s error:%s", image.Url, err)
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

func (t *TwitterClient) uploadVideos(item *crawler.ChannelItem) []*anaconda.VideoMedia {
	videos := []*anaconda.VideoMedia{}
	for _, video := range item.Videos {
		resp, err := t.Api.HttpClient.Get(video.Url)
		if err != nil {
			t.Log.Errorf("failed to get mp4 url:%s err:%v\n", video.Url, err)
			continue
		}

		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Log.Errorf("failed to readall mp4 url:%s err:%v\n", video.Url, err)
			continue
		}

		media, err := t.Api.UploadVideoInit(len(bytes), "video/mp4")
		if err != nil {
			t.Log.Errorf("failed to upload video init. url:%s err:%v\n", video.Url, err)
			continue
		}
		if err = t.Api.UploadVideoAppend(media.MediaIDString, 0, base64.StdEncoding.EncodeToString(bytes)); err != nil {
			t.Log.Errorf("failed to upload video append. url:%s err:%v\n", video.Url, err)
			continue
		}
		v, err := t.Api.UploadVideoFinalize(media.MediaIDString)
		if err != nil {
			t.Log.Errorf("failed to upload video finalize. url:%s err:%v\n", video.Url, err)
			continue
		}
		videos = append(videos, &v)
	}
	return videos
}

func (t *TwitterClient) downloadImage(url string) (string, error) {
	response, err := t.Api.HttpClient.Get(url)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to download image. url:%s", url)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to read response. url:%s", url)
	}
	return base64.StdEncoding.EncodeToString(body), nil
}
