package model

type (
	// LineNotifyToken interface
	LineNotifyToken interface {
		GetAccessToken(string) (string, error)
	}

	// LineNotifyRequest represents request that line notify message, img urls
	LineNotifyRequest struct {
		Text     string
		ImageURL string
	}

	// LineNotifyResponse represents response line notify data
	LineNotifyResponse struct {
	}

	// LineNotify interface
	LineNotify interface {
		Notify(LineNotifyRequest) (LineNotifyResponse, error)
	}
)
