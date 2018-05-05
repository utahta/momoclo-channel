package usecase

import (
	"context"

	"github.com/utahta/momoclo-channel/customsearch"
	"github.com/utahta/momoclo-channel/linebot"
	"github.com/utahta/momoclo-channel/log"
)

type (
	// HandleLineBotEvents use case
	HandleLineBotEvents struct {
		log           log.Logger
		lineBot       linebot.Client
		imageSearcher customsearch.ImageSearcher
	}

	// HandleLineBotEventsParams use case params
	HandleLineBotEventsParams struct {
		Events []linebot.Event
	}
)

// NewHandleLineBotEvents returns HandleLineBotEvents use case
func NewHandleLineBotEvents(
	logger log.Logger,
	lineBot linebot.Client,
	imageSearcher customsearch.ImageSearcher) *HandleLineBotEvents {
	return &HandleLineBotEvents{
		log:           logger,
		lineBot:       lineBot,
		imageSearcher: imageSearcher,
	}
}

// Do handles given line bot events
func (use *HandleLineBotEvents) Do(ctx context.Context, params HandleLineBotEventsParams) error {
	const errTag = "HandleLineBotEvents.Do"

	for _, event := range params.Events {
		use.log.Infof("handle event:%v", event)

		switch event.Type {
		case linebot.EventTypeMessage:
			switch event.MessageType {
			case linebot.MessageTypeText:
				if linebot.MatchOn(event.TextMessage.Text) {
					use.lineBot.ReplyText(ctx, event.ReplyToken, linebot.OnMessage())
					continue
				} else if linebot.MatchOff(event.TextMessage.Text) {
					use.lineBot.ReplyText(ctx, event.ReplyToken, linebot.OffMessage())
					continue
				}

				memberName := linebot.FindMemberName(event.TextMessage.Text)
				if memberName == "" {
					use.lineBot.ReplyText(ctx, event.ReplyToken, linebot.HelpMessage())
					continue
				}

				img, err := use.imageSearcher.Search(memberName)
				if err != nil {
					use.log.Warningf("%v: image not found word:%v err:%v", errTag, memberName, err)
					use.lineBot.ReplyText(ctx, event.ReplyToken, linebot.ImageNotFoundMessage())
					continue
				}
				use.lineBot.ReplyImage(ctx, event.ReplyToken, img.URL, img.ThumbnailURL)

			default:
				use.log.Infof("not handle message type:%v", event.MessageType)
			}
		case linebot.EventTypeFollow:
			use.log.Info("follow event")
			use.lineBot.ReplyText(ctx, event.ReplyToken, linebot.FollowMessage())
		case linebot.EventTypeUnfollow:
			use.log.Info("unfollow event")
		default:
			use.log.Info("not handle event type:%v", event.Type)
		}
	}
	return nil
}
