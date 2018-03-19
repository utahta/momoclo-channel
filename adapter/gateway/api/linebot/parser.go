package linebot

import (
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/utahta/momoclo-channel/config"
	"github.com/utahta/momoclo-channel/domain"
	"github.com/utahta/momoclo-channel/domain/model"
)

// ParseRequest parses http request
func ParseRequest(r *http.Request) ([]model.LineBotEvent, error) {
	events, err := linebot.ParseRequest(config.C.LineBot.ChannelSecret, r)
	if err == linebot.ErrInvalidSignature {
		return nil, domain.ErrInvalidSignature
	} else if err != nil {
		return nil, err
	}

	results := make([]model.LineBotEvent, len(events))
	for i, event := range events {
		results[i].ReplyToken = event.ReplyToken

		switch event.Type {
		case linebot.EventTypeMessage:
			results[i].Type = model.LineBotEventTypeMessage
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				results[i].TextMessage = model.LineBotTextMessage{ID: message.ID, Text: message.Text}
				results[i].MessageType = model.LineBotMessageTypeText
			}
		case linebot.EventTypeFollow:
			results[i].Type = model.LineBotEventTypeFollow
		case linebot.EventTypeUnfollow:
			results[i].Type = model.LineBotEventTypeUnfollow
		}
	}
	return results, nil
}
