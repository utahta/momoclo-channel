package usecase

import (
	"github.com/utahta/momoclo-channel/domain/core"
	"github.com/utahta/momoclo-channel/domain/model"
	"github.com/utahta/momoclo-channel/domain/service/linebot"
)

type (
	// HandleLineBotEvents use case
	HandleLineBotEvents struct {
		log          core.Logger
		lineBot      model.LineBotClient
		customSearch model.CustomSearchClient
	}

	// HandleLineBotEventsParams use case params
	HandleLineBotEventsParams struct {
		Events []model.LineBotEvent
	}
)

// NewHandleLineBotEvents returns HandleLineBotEvents use case
func NewHandleLineBotEvents(
	logger core.Logger,
	lineBot model.LineBotClient,
	customSearch model.CustomSearchClient) *HandleLineBotEvents {
	return &HandleLineBotEvents{
		log:          logger,
		lineBot:      lineBot,
		customSearch: customSearch,
	}
}

// Do handles given line bot events
func (use *HandleLineBotEvents) Do(params HandleLineBotEventsParams) error {
	const errTag = "HandleLineBotEvents.Do"

	for _, event := range params.Events {
		use.log.Infof("handle event:%v", event)

		switch event.Type {
		case model.LineBotEventTypeMessage:
			switch event.MessageType {
			case model.LineBotMessageTypeText:
				if linebot.MatchOn(event.TextMessage.Text) {
					use.lineBot.ReplyText(event.ReplyToken, linebot.OnMessage())
					continue
				} else if linebot.MatchOff(event.TextMessage.Text) {
					use.lineBot.ReplyText(event.ReplyToken, linebot.OffMessage())
					continue
				}

				memberName := linebot.FindMemberName(event.TextMessage.Text)
				if memberName == "" {
					use.lineBot.ReplyText(event.ReplyToken, linebot.HelpMessage())
					continue
				}

				img, err := use.customSearch.SearchImage(memberName)
				if err != nil {
					use.log.Warningf("%v: image not found word:%v", err, memberName)
					use.lineBot.ReplyText(event.ReplyToken, linebot.ImageNotFoundMessage())
					continue
				}
				use.lineBot.ReplyImage(event.ReplyToken, img.URL, img.ThumbnailURL)

			default:
				use.log.Infof("not handle message type:%v", event.MessageType)
			}
		case model.LineBotEventTypeFollow:
			use.log.Info("follow event")
			use.lineBot.ReplyText(event.ReplyToken, linebot.FollowMessage())
		case model.LineBotEventTypeUnfollow:
			use.log.Info("unfollow event")
		default:
			use.log.Info("not handle event type:%v", event.Type)
		}
	}
	return nil
}
