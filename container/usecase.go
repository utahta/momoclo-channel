package container

import (
	"context"

	"github.com/utahta/momoclo-channel/adapter/gateway/api/twitter"
	"github.com/utahta/momoclo-channel/adapter/gateway/crawler"
	"github.com/utahta/momoclo-channel/infrastructure/event"
	"github.com/utahta/momoclo-channel/infrastructure/log"
	"github.com/utahta/momoclo-channel/usecase"
)

// UsecaseContainer dependency injection
type UsecaseContainer struct {
	ctx  context.Context
	repo *RepositoryContainer
}

// Usecase returns container of use case
func Usecase(ctx context.Context) *UsecaseContainer {
	return &UsecaseContainer{ctx, Repository(ctx)}
}

// FeedsCrawl use case
func (c *UsecaseContainer) FeedsCrawl() *usecase.FeedsCrawl {
	return usecase.NewFeedsCrawl(
		c.ctx,
		log.NewAppengineLogger(c.ctx),
		c.FeedCrawl(),
	)
}

// FeedCrawl use case
func (c *UsecaseContainer) FeedCrawl() *usecase.FeedCrawl {
	return usecase.NewFeedCrawl(
		c.ctx,
		log.NewAppengineLogger(c.ctx),
		crawler.New(c.ctx),
		event.NewTaskQueue(c.ctx),
		c.repo.LatestEntryRepository(),
	)
}

// FeedTweetsEnqueue use case
func (c *UsecaseContainer) FeedTweetsEnqueue() *usecase.FeedTweetsEnqueue {
	return usecase.NewFeedTweetsEnqueue(
		log.NewAppengineLogger(c.ctx),
		event.NewTaskQueue(c.ctx),
	)
}

// FeedTweet use case
func (c *UsecaseContainer) FeedTweet() *usecase.FeedTweet {
	logger := log.NewAppengineLogger(c.ctx)
	return usecase.NewFeedTweet(
		logger,
		event.NewTaskQueue(c.ctx),
		twitter.NewTweeter(c.ctx, logger),
	)
}
