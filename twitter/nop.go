package twitter

type nop struct{}

// NewNopTweeter returns no operation tweeter
func NewNopTweeter() Tweeter {
	return &nop{}
}

func (c *nop) Tweet(req TweetRequest) (TweetResponse, error) {
	return TweetResponse{}, nil
}
