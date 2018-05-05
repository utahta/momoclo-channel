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
		GetAccessToken(context.Context, string) (string, error)
	}

	tokenClient struct {
	}
)

// NewToken returns Token
func NewToken() Token {
	return &tokenClient{}
}

// GetAccessToken returns access token that published by LINE Notify
func (c *tokenClient) GetAccessToken(ctx context.Context, code string) (string, error) {
	t := token.New(
		CallbackURL(),
		config.C().LineNotify.ClientID,
		config.C().LineNotify.ClientSecret,
		token.WithHTTPClient(urlfetch.Client(ctx)),
	)
	return t.GetAccessToken(code)
}
