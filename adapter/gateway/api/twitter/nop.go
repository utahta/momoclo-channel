package twitter

import (
	"github.com/utahta/momoclo-channel/types"
)

type nop struct{}

// NewNopTweeter returns no operation tweeter
func NewNopTweeter() types.Tweeter {
	return &nop{}
}

func (c *nop) Tweet(req types.TweetRequest) (types.TweetResponse, error) {
	return types.TweetResponse{}, nil
}
