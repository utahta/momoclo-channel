package container

import (
	"context"

	"github.com/utahta/momoclo-channel/usecase"
)

type usecaseContainer struct {
	ctx  context.Context
	repo *repositoryContainer
}

// Usecase returns use case container
func Usecase(ctx context.Context) *usecaseContainer {
	return &usecaseContainer{ctx, Repository(ctx)}
}

func (c *usecaseContainer) Crawl() *usecase.Crawl {
	return usecase.NewCrawl(c.repo.LatestEntryRepository())
}
