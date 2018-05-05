package linenotify

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/utahta/go-linenotify"
	"github.com/utahta/go-openuri"
	"github.com/utahta/momoclo-channel/config"
	"github.com/utahta/nsync"
	"google.golang.org/appengine/urlfetch"
)

type (
	// Message represents text message and image
	Message struct {
		Text     string `validate:"required"`
		ImageURL string `validate:"omitempty,url"`
	}

	// Request represents request that notification message
	Request struct {
		ID          string    `validate:"required"`
		AccessToken string    `validate:"required"`
		Messages    []Message `validate:"min=1,dive"`
	}

	// Client interface
	Client interface {
		Notify(context.Context, string, Message) error
	}

	client struct {
	}
)

var (
	cacheRepo     = newCacheRepository()
	cacheNamedMux nsync.Mutex
)

// New returns LineNotify
func New() Client {
	if config.C().LineNotify.Disabled {
		return NewNop()
	}
	return &client{}
}

// Notify sends message to given token
func (c *client) Notify(ctx context.Context, accessToken string, msg Message) error {
	if err := c.notify(ctx, accessToken, msg); err != nil {
		if err == linenotify.ErrNotifyInvalidAccessToken {
			return ErrInvalidAccessToken
		}
		return err
	}
	return nil
}

func (c *client) notify(ctx context.Context, accessToken string, msg Message) error {
	notify := linenotify.New()
	notify.HTTPClient = urlfetch.Client(ctx)

	if msg.ImageURL != "" {
		b, err := c.fetchImage(ctx, msg.ImageURL)
		if err != nil {
			return err
		}
		if _, err := notify.NotifyWithImage(accessToken, msg.Text, bytes.NewReader(b)); err != nil {
			return err
		}
	} else {
		if _, err := notify.NotifyMessage(accessToken, msg.Text); err != nil {
			return err
		}
	}
	return nil
}

func (c *client) fetchImage(ctx context.Context, urlStr string) ([]byte, error) {
	cacheNamedMux.Lock(urlStr)
	defer cacheNamedMux.Unlock(urlStr)

	if c, ok := cacheRepo.Get(urlStr); ok {
		return c.Bytes(), nil
	}

	o, err := openuri.Open(urlStr, openuri.WithHTTPClient(urlfetch.Client(ctx)))
	if err != nil {
		return nil, err
	}
	defer o.Close()

	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, o); err != nil {
		return nil, err
	}

	if ct := http.DetectContentType(buf.Bytes()); !strings.Contains(ct, "image") {
		return nil, errors.Errorf("invalid content type. ct:%v", ct)
	}
	cacheRepo.Set(urlStr, newCache(buf.Bytes()))

	return buf.Bytes(), nil
}
