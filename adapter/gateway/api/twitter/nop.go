package twitter

import (
	"github.com/utahta/momoclo-channel/domain/model"
)

type nop struct{}

// NewNopTweeter returns no operation tweeter
func NewNopTweeter() model.Tweeter {
	return &nop{}
}

func (c *nop) Tweet(req model.TweetRequest) (model.TweetResponse, error) {
	return model.TweetResponse{}, nil
}
