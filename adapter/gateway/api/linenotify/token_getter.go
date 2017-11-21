package linenotify

import (
	"context"

	"github.com/utahta/go-linenotify"
	"github.com/utahta/momoclo-channel/domain/model"
	"github.com/utahta/momoclo-channel/lib/config"
	"google.golang.org/appengine/urlfetch"
)

type (
	tokenGetter struct {
		*linenotify.TokenClient
	}
)

// NewTokenGetter returns LineNotifyTokenClient
func NewTokenGetter(ctx context.Context) model.LineNotifyTokenGetter {
	c := &tokenGetter{
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
func (c *tokenGetter) Get(code string) (string, error) {
	c.TokenClient.Code = code
	return c.TokenClient.Get()
}
