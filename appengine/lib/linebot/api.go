package linebot

import (
	"github.com/utahta/momoclo-channel/appengine/lib/log"
	"github.com/utahta/momoclo-channel/crawler"
	"golang.org/x/net/context"
)

func NotifyChannel(ctx context.Context, title string, item *crawler.ChannelItem) {
	bot, err := NewClient(ctx)
	if err != nil {
		log.GaeLog(ctx).Error(err)
		return
	}

	if err := bot.NotifyChannel(title, item); err != nil {
		log.GaeLog(ctx).Error(err)
		return
	}
}

func NotifyMessage(ctx context.Context, text string) {
	bot, err := NewClient(ctx)
	if err != nil {
		log.GaeLog(ctx).Error(err)
		return
	}

	if err := bot.NotifyMessage(text); err != nil {
		log.GaeLog(ctx).Error(err)
		return
	}
}

func NotifyMessageTo(ctx context.Context, to []string, text string) {
	bot, err := NewClient(ctx)
	if err != nil {
		log.GaeLog(ctx).Error(err)
		return
	}

	if err := bot.NotifyMessageTo(to, text); err != nil {
		log.GaeLog(ctx).Error(err)
		return
	}
}
