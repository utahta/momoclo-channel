package model

type (
	// LineNotifyToken interface
	LineNotifyToken interface {
		GetAccessToken(string) (string, error)
	}

	// LineNotifyMessage represents text message and image
	LineNotifyMessage struct {
		Text     string `validate:"required"`
		ImageURL string `validate:"omitempty,url"`
	}

	// LineNotifyRequest represents request that notification message
	LineNotifyRequest struct {
		ID          string              `validate:"required"`
		AccessToken string              `validate:"required"`
		Messages    []LineNotifyMessage `validate:"min=1,dive"`
	}

	// LineNotify interface
	LineNotify interface {
		Notify(string, LineNotifyMessage) error
	}
)
