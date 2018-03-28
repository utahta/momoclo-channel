package linebot

import (
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/utahta/momoclo-channel/config"
	"github.com/utahta/momoclo-channel/types"
)

// ParseRequest parses http request
func ParseRequest(r *http.Request) ([]types.LineBotEvent, error) {
	events, err := linebot.ParseRequest(config.C.LineBot.ChannelSecret, r)
	if err == linebot.ErrInvalidSignature {
		return nil, types.ErrInvalidSignature
	} else if err != nil {
		return nil, err
	}

	results := make([]types.LineBotEvent, len(events))
	for i, event := range events {
		results[i].ReplyToken = event.ReplyToken

		switch event.Type {
		case linebot.EventTypeMessage:
			results[i].Type = types.LineBotEventTypeMessage
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				results[i].TextMessage = types.LineBotTextMessage{ID: message.ID, Text: message.Text}
				results[i].MessageType = types.LineBotMessageTypeText
			}
		case linebot.EventTypeFollow:
			results[i].Type = types.LineBotEventTypeFollow
		case linebot.EventTypeUnfollow:
			results[i].Type = types.LineBotEventTypeUnfollow
		}
	}
	return results, nil
}
