package linebot

import (
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/utahta/momoclo-channel/config"
)

type (
	// EventType represents line bot event type
	EventType string

	// MessageType represents line bot message type
	MessageType string

	// TextMessage represents line bot text message
	TextMessage struct {
		ID   string
		Text string
	}

	// Event represents line bot event
	Event struct {
		ReplyToken  string
		Type        EventType
		MessageType MessageType
		TextMessage TextMessage
	}
)

const (
	EventTypeMessage  EventType = "message"
	EventTypeFollow   EventType = "follow"
	EventTypeUnfollow EventType = "unfollow"

	MessageTypeText MessageType = "text"
)

// ParseRequest parses http request
func ParseRequest(r *http.Request) ([]Event, error) {
	events, err := linebot.ParseRequest(config.C.LineBot.ChannelSecret, r)
	if err == linebot.ErrInvalidSignature {
		return nil, ErrInvalidSignature
	} else if err != nil {
		return nil, err
	}

	results := make([]Event, len(events))
	for i, event := range events {
		results[i].ReplyToken = event.ReplyToken

		switch event.Type {
		case linebot.EventTypeMessage:
			results[i].Type = EventTypeMessage
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				results[i].TextMessage = TextMessage{ID: message.ID, Text: message.Text}
				results[i].MessageType = MessageTypeText
			}
		case linebot.EventTypeFollow:
			results[i].Type = EventTypeFollow
		case linebot.EventTypeUnfollow:
			results[i].Type = EventTypeUnfollow
		}
	}
	return results, nil
}
