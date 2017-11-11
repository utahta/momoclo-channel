package container

import (
	"context"

	"github.com/utahta/momoclo-channel/adapter/persistence"
	"github.com/utahta/momoclo-channel/domain/entity"
	"github.com/utahta/momoclo-channel/infrastructure/dao"
)

type repositoryContainer struct {
	ctx context.Context
}

// Repository returns repositories container
func Repository(ctx context.Context) *repositoryContainer {
	return &repositoryContainer{ctx}
}

func (c *repositoryContainer) LatestEntryRepository() entity.LatestEntryRepository {
	return persistence.NewLatestEntryRepository(dao.NewDatastoreHandler(c.ctx))
}
