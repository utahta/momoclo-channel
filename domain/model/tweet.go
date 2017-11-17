package model

type (
	// TweetRequest represents request that tweet message, img urls and video url data
	TweetRequest struct {
		InReplyToStatusID string
		Text              string
		ImageURLs         []string
		VideoURL          string
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
