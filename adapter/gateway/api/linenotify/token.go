package linenotify

import (
	"context"

	"github.com/utahta/go-linenotify/token"
	"github.com/utahta/momoclo-channel/domain/model"
	"github.com/utahta/momoclo-channel/lib/config"
	"google.golang.org/appengine/urlfetch"
)

type (
	tokenClient struct {
		*token.Client
	}
)

// NewToken returns LineNotifyToken
func NewToken(ctx context.Context) model.LineNotifyToken {
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
