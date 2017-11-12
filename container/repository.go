package container

import (
	"context"

	"github.com/utahta/momoclo-channel/adapter/persistence"
	"github.com/utahta/momoclo-channel/domain/entity"
	"github.com/utahta/momoclo-channel/infrastructure/dao"
)

// RepositoryContainer dependency injection
type RepositoryContainer struct {
	ctx context.Context
}

// Repository returns repositories container
func Repository(ctx context.Context) *RepositoryContainer {
	return &RepositoryContainer{ctx}
}

// LatestEntryRepository returns LatestEntry repository
func (c *RepositoryContainer) LatestEntryRepository() entity.LatestEntryRepository {
	return persistence.NewLatestEntryRepository(dao.NewDatastoreHandler(c.ctx))
}
