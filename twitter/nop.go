package twitter

import "context"

type nop struct{}

// NewNopTweeter returns no operation tweeter
func NewNopTweeter() Tweeter {
	return &nop{}
}

func (c *nop) Tweet(_ context.Context, _ TweetRequest) (TweetResponse, error) {
	return TweetResponse{}, nil
}
