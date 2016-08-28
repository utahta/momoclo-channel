package twitter

import (
	"net/url"

	"github.com/ChimeraCoder/anaconda"
	"github.com/pkg/errors"
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

func (t *MessageClient) Tweet(msg string) error {
	v := url.Values{}
	_, err := t.Api.PostTweet(msg, v)
	if err != nil {
		return errors.Wrapf(err, "Failed to post message. msg:%s", msg)
	}
	t.Log.Infof("Post message. msg:%s", msg)
	return nil
}
