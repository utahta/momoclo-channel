package linebot

import (
	"github.com/utahta/momoclo-channel/appengine/lib/log"
	"github.com/utahta/momoclo-channel/crawler"
	"golang.org/x/net/context"
)

func NotifyChannel(ctx context.Context, title string, item *crawler.ChannelItem) {
	bot, err := Dial(ctx)
	if err != nil {
		log.GaeLog(ctx).Error(err)
		return
	}
	defer bot.Close()

	if err := bot.NotifyChannel(title, item); err != nil {
		log.GaeLog(ctx).Error(err)
		return
	}
}

func NotifyUstream(ctx context.Context) {
	bot, err := Dial(ctx)
	if err != nil {
		log.GaeLog(ctx).Error(err)
		return
	}
	defer bot.Close()

	if err := bot.NotifyUstream(); err != nil {
		log.GaeLog(ctx).Error(err)
		return
	}
}

func NotifyReminder(ctx context.Context, text string) {
	bot, err := Dial(ctx)
	if err != nil {
		log.GaeLog(ctx).Error(err)
		return
	}
	defer bot.Close()

	if err := bot.NotifyReminder(text); err != nil {
		log.GaeLog(ctx).Error(err)
		return
	}
}