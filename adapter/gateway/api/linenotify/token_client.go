package linenotify

import (
	"context"

	"github.com/utahta/go-linenotify"
	"github.com/utahta/momoclo-channel/domain/model"
	"github.com/utahta/momoclo-channel/lib/config"
	"google.golang.org/appengine/urlfetch"
)

type (
	tokenClient struct {
		*linenotify.TokenClient
	}
)

// NewTokenClient returns LineNotifyTokenClient
func NewTokenClient(ctx context.Context) model.LineNotifyTokenClient {
	c := &tokenClient{
		linenotify.NewToken(
			"",
			config.LineNotifyCallbackURL(),
			config.C.LineNotify.ClientID,
			config.C.LineNotify.ClientSecret,
		),
	}
	c.HTTPClient = urlfetch.Client(ctx)
	return c
}

// GetToken returns LINE Notify token
func (c *tokenClient) GetToken(code string) (string, error) {
	c.TokenClient.Code = code
	return c.TokenClient.Get()
}
