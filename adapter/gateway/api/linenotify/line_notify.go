package linenotify

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/utahta/go-linenotify"
	"github.com/utahta/go-openuri"
	"github.com/utahta/momoclo-channel/domain"
	"github.com/utahta/momoclo-channel/domain/model"
	"github.com/utahta/momoclo-channel/lib/config"
	"google.golang.org/appengine/urlfetch"
)

type (
	client struct {
		*linenotify.Client
	}

	cacheRepository map[string][]byte
)

var (
	cachePool = sync.Pool{
		New: func() interface{} {
			return cacheRepository{}
		},
	}
	cacheMux sync.Mutex
)

// New returns LineNotify
func New(ctx context.Context) model.LineNotify {
	if config.C.LineNotify.Disabled {
		return NewNop()
	}

	c := linenotify.New()
	c.HTTPClient = urlfetch.Client(ctx)
	return &client{Client: c}
}

// Notify sends message to given token
func (c *client) Notify(accessToken string, msg model.LineNotifyMessage) error {
	if err := c.notify(accessToken, msg); err != nil {
		if err == linenotify.ErrNotifyInvalidAccessToken {
			return domain.ErrInvalidAccessToken
		}
		return err
	}
	return nil
}

func (c *client) notify(accessToken string, msg model.LineNotifyMessage) error {
	if msg.ImageURL != "" {
		b, err := c.fetchImage(msg.ImageURL)
		if err != nil {
			return err
		}
		if _, err := c.Client.NotifyWithImage(accessToken, msg.Text, bytes.NewReader(b)); err != nil {
			return err
		}
	} else {
		if _, err := c.Client.NotifyMessage(accessToken, msg.Text); err != nil {
			return err
		}
	}
	return nil
}

func (c *client) fetchImage(urlStr string) ([]byte, error) {
	cacheMux.Lock()
	defer cacheMux.Unlock()

	cache := cachePool.Get().(cacheRepository)
	defer cachePool.Put(cache)

	if b, ok := cache[urlStr]; ok {
		return b, nil
	}

	o, err := openuri.Open(urlStr, openuri.WithHTTPClient(c.HTTPClient))
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
	cache[urlStr] = buf.Bytes()

	return buf.Bytes(), nil
}
