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

// CrawlFeeds use case
func (c *UsecaseContainer) CrawlFeeds() *usecase.CrawlFeeds {
	return usecase.NewCrawlFeeds(
		log.NewAppengineLogger(c.ctx),
		c.CrawlFeed(),
	)
}

// CrawlFeed use case
func (c *UsecaseContainer) CrawlFeed() *usecase.CrawlFeed {
	return usecase.NewCrawlFeed(
		log.NewAppengineLogger(c.ctx),
		crawler.New(c.ctx),
		event.NewTaskQueue(c.ctx),
		c.repo.LatestEntryRepository(),
	)
}

// EnqueueTweets use case
func (c *UsecaseContainer) EnqueueTweets() *usecase.EnqueueTweets {
	return usecase.NewEnqueueTweets(
		log.NewAppengineLogger(c.ctx),
		event.NewTaskQueue(c.ctx),
	)
}

// Tweet use case
func (c *UsecaseContainer) Tweet() *usecase.Tweet {
	return usecase.NewTweet(
		log.NewAppengineLogger(c.ctx),
		event.NewTaskQueue(c.ctx),
		twitter.NewTweeter(c.ctx),
	)
}
