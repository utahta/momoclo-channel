package container

import (
	"context"

	"github.com/utahta/momoclo-channel/adapter/gateway/persistence"
	"github.com/utahta/momoclo-channel/dao"
	"github.com/utahta/momoclo-channel/entity"
	"github.com/utahta/momoclo-channel/types"
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
func (c *RepositoryContainer) LatestEntryRepository() types.LatestEntryRepository {
	return persistence.NewLatestEntryRepository(dao.NewDatastoreHandler(c.ctx))
}

// TweetItemRepository returns TweetItem repository
func (c *RepositoryContainer) TweetItemRepository() types.TweetItemRepository {
	return persistence.NewTweetItemRepository(dao.NewDatastoreHandler(c.ctx))
}

// LineItemRepository returns LineItem repository
func (c *RepositoryContainer) LineItemRepository() types.LineItemRepository {
	return persistence.NewLineItemRepository(dao.NewDatastoreHandler(c.ctx))
}

// ReminderRepository returns Reminder repository
func (c *RepositoryContainer) ReminderRepository() types.ReminderRepository {
	return persistence.NewReminderRepository(dao.NewDatastoreHandler(c.ctx))
}

// UstreamStatusRepository returns UstreamStatus repository
func (c *RepositoryContainer) UstreamStatusRepository() entity.UstreamStatusRepository {
	return entity.NewUstreamStatusRepository(dao.NewDatastoreHandler(c.ctx))
}

// LineNotificationRepository returns LineNotification repository
func (c *RepositoryContainer) LineNotificationRepository() types.LineNotificationRepository {
	return persistence.NewLineNotificationRepository(dao.NewDatastoreHandler(c.ctx))
}
