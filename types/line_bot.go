package types

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

	// LineBot represents line bot conversations
	LineBot interface {
		ReplyText(string, string) error
		ReplyImage(string, string, string) error
	}
)

const (
	LineBotEventTypeMessage  LineBotEventType = "message"
	LineBotEventTypeFollow   LineBotEventType = "follow"
	LineBotEventTypeUnfollow LineBotEventType = "unfollow"

	LineBotMessageTypeText LineBotMessageType = "text"
)
