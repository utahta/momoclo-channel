package linenotify

import (
	"context"
	"fmt"
	"time"

	"github.com/utahta/momoclo-crawler"
)

const timeout = 540 * time.Second

type ChannelParam struct {
	Title string
	Item  *crawler.ChannelItem
}

// Send message to LINE Notify
func NotifyMessage(ctx context.Context, message string) error {
	reqCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	c, err := newClient(reqCtx)
	if err != nil {
		return err
	}

	// [Notify Name] が付くので先頭に改行をいれて調整
	return c.notifyMessage(fmt.Sprintf("\n%s", message), "")
}

// Send channel message and images to LINE Notify
func NotifyChannel(ctx context.Context, param *ChannelParam) error {
	reqCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	c, err := newClient(reqCtx)
	if err != nil {
		return err
	}

	return c.notifyChannelItem(param)
}
