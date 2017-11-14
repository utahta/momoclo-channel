package container

import (
	"context"

	"github.com/utahta/momoclo-channel/adapter/gateway/persistence"
	"github.com/utahta/momoclo-channel/domain/model"
	"github.com/utahta/momoclo-channel/infrastructure/dao"
)

// RepositoryContainer dependency injection
type RepositoryContainer struct {
	ctx context.Context
}

// Repository returns container of repositories
func Repository(ctx context.Context) *RepositoryContainer {
	return &RepositoryContainer{ctx}
}

// LatestEntryRepository returns LatestEntry repository
func (c *RepositoryContainer) LatestEntryRepository() model.LatestEntryRepository {
	return persistence.NewLatestEntryRepository(dao.NewDatastoreHandler(c.ctx))
}
