package linenotify

import (
	"context"

	"github.com/utahta/go-linenotify/token"
	"github.com/utahta/momoclo-channel/config"
	"google.golang.org/appengine/urlfetch"
)

type (
	// Token interface
	Token interface {
		GetAccessToken(string) (string, error)
	}

	tokenClient struct {
		*token.Client
	}
)

// NewToken returns Token
func NewToken(ctx context.Context) Token {
	return &tokenClient{
		token.New(
			config.LineNotifyCallbackURL(),
			config.C.LineNotify.ClientID,
			config.C.LineNotify.ClientSecret,
			token.WithHTTPClient(urlfetch.Client(ctx)),
		),
	}
}

// GetAccessToken returns access token that published by LINE Notify
func (c *tokenClient) GetAccessToken(code string) (string, error) {
	return c.Client.GetAccessToken(code)
}
