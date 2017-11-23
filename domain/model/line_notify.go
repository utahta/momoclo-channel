package model

type (
	// LineNotifyToken interface
	LineNotifyToken interface {
		GetAccessToken(string) (string, error)
	}

	// LineNotifyMessage represents text message and image
	LineNotifyMessage struct {
		Text     string
		ImageURL string
	}

	// LineNotifyRequest represents request that notification message
	LineNotifyRequest struct {
		AccessToken string
		Messages    []LineNotifyMessage
	}

	// LineNotify interface
	LineNotify interface {
		Notify(string, LineNotifyMessage) error
	}
)
