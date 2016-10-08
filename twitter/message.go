package twitter

import (
	"net/url"

	"github.com/pkg/errors"
)

type MessageClient struct {
	*Client
}

func NewMessageClient(consumerKey, consumerSecret, accessToken, accessTokenSecret string, options ...ClientOption) (*MessageClient, error) {
	c, err := newClient(consumerKey, consumerSecret, accessToken, accessTokenSecret, options...)
	if err != nil {
		return nil, err
	}
	return &MessageClient{Client: c}, nil
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
