package container

import (
	"context"

	"github.com/utahta/momoclo-channel/usecase"
)

// UsecaseContainer dependency injection
type UsecaseContainer struct {
	ctx  context.Context
	repo *RepositoryContainer
}

// Usecase returns use case container
func Usecase(ctx context.Context) *UsecaseContainer {
	return &UsecaseContainer{ctx, Repository(ctx)}
}

// Crawl returns Crawl use case
func (c *UsecaseContainer) Crawl() *usecase.Crawl {
	return usecase.NewCrawl(c.repo.LatestEntryRepository())
}
