package twitter

import (
	"net/url"

	"github.com/ChimeraCoder/anaconda"
	"github.com/utahta/momoclo-channel/log"
)

type MessageClient struct {
	Api *anaconda.TwitterApi
	Log log.Logger
}

func NewMessageClient(consumerKey, consumerSecret, accessToken, accessTokenSecret string) *MessageClient {
	t := &MessageClient{}
	t.auth(consumerKey, consumerSecret, accessToken, accessTokenSecret)
	t.Log = log.NewSilentLogger()
	return t
}

func (t *MessageClient) auth(consumerKey, consumerSecret, accessToken, accessTokenSecret string) {
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	t.Api = anaconda.NewTwitterApi(accessToken, accessTokenSecret)
}

func (t *MessageClient) Tweet(msg string) {
	v := url.Values{}
	_, err := t.Api.PostTweet(msg, v)
	if err != nil {
		t.Log.Errorf("Failed to post message. msg:%s error:%s", msg, err)
		return
	}
	t.Log.Infof("Post message. msg:%s", msg)
}
