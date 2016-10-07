package twitter

import (
	"net/http"

	"github.com/ChimeraCoder/anaconda"
	"github.com/utahta/momoclo-channel/log"
)

type Client struct {
	Api *anaconda.TwitterApi
	Log log.Logger
}

type ClientOption func(*Client) error

func newClient(consumerKey, consumerSecret, accessToken, accessTokenSecret string, options ...ClientOption) (*Client, error) {
	c := &Client{}
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	c.Api = anaconda.NewTwitterApi(accessToken, accessTokenSecret)

	for _, option := range options {
		if err := option(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// WithHTTPTransport function
func WithHTTPTransport(t http.RoundTripper) ClientOption {
	return func(client *Client) error {
		client.Api.HttpClient.Transport = t
		return nil
	}
}

// WithLogger function
func WithLogger(l log.Logger) ClientOption {
	return func(client *Client) error {
		client.Log = l
		return nil
	}
}
