package model

type (
	// LineBotEventType represents line bot event type
	LineBotEventType string

	// LineBotMessageType represents line bot message type
	LineBotMessageType string

	// LineBotTextMessage represents line bot text message
	LineBotTextMessage struct {
		ID   string
		Text string
	}

	// LineBotEvent represents line bot event
	LineBotEvent struct {
		ReplyToken  string
		Type        LineBotEventType
		MessageType LineBotMessageType
		TextMessage LineBotTextMessage
	}

	// LineBotClient represents line bot dispatch
	LineBotClient interface {
		ReplyText(string, string) error
		ReplyImage(string, string, string) error
	}
)

const (
	LineBotEventTypeNone     LineBotEventType = ""
	LineBotEventTypeMessage  LineBotEventType = "message"
	LineBotEventTypeFollow   LineBotEventType = "follow"
	LineBotEventTypeUnfollow LineBotEventType = "unfollow"

	LineBotMessageTypeNone LineBotMessageType = ""
	LineBotMessageTypeText LineBotMessageType = "text"
)
