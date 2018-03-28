package types

type (
	// TweetRequest represents request that tweet message, img urls and video url data
	TweetRequest struct {
		InReplyToStatusID string
		Text              string
		ImageURLs         []string `validate:"dive,omitempty,url"`
		VideoURL          string   `validate:"omitempty,url"`
	}

	// TweetResponse represents response tweet data
	TweetResponse struct {
		IDStr string
	}

	// Tweeter interface
	Tweeter interface {
		Tweet(TweetRequest) (TweetResponse, error)
	}
)
