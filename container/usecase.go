package container

import (
	"context"

	"github.com/utahta/momoclo-channel/infrastructure/crawler"
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

// CrawlAll returns CrawlAll use case
func (c *UsecaseContainer) CrawlAll() *usecase.CrawlAll {
	return usecase.NewCrawlAll(
		c.ctx,
		log.NewAppengineLogger(c.ctx),
		c.Crawl(),
	)
}

// Crawl returns Crawl use case
func (c *UsecaseContainer) Crawl() *usecase.Crawl {
	return usecase.NewCrawl(
		c.ctx,
		log.NewAppengineLogger(c.ctx),
		crawler.New(c.ctx),
		event.NewTaskQueue(c.ctx),
		c.repo.LatestEntryRepository(),
	)
}
