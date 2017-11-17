package twitter

import (
	"github.com/utahta/momoclo-channel/domain/model"
)

type nop struct{}

func (c *nop) Tweet(req model.TweetRequest) (model.TweetResponse, error) {
	return model.TweetResponse{}, nil
}
