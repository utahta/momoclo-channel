package container

import (
	"context"

	"github.com/utahta/momoclo-channel/dao"
	"github.com/utahta/momoclo-channel/entity"
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
func (c *RepositoryContainer) LatestEntryRepository() entity.LatestEntryRepository {
	return entity.NewLatestEntryRepository(dao.NewDatastoreHandler(c.ctx))
}

// TweetItemRepository returns TweetItem repository
func (c *RepositoryContainer) TweetItemRepository() entity.TweetItemRepository {
	return entity.NewTweetItemRepository(dao.NewDatastoreHandler(c.ctx))
}

// LineItemRepository returns LineItem repository
func (c *RepositoryContainer) LineItemRepository() entity.LineItemRepository {
	return entity.NewLineItemRepository(dao.NewDatastoreHandler(c.ctx))
}

// ReminderRepository returns Reminder repository
func (c *RepositoryContainer) ReminderRepository() entity.ReminderRepository {
	return entity.NewReminderRepository(dao.NewDatastoreHandler(c.ctx))
}

// UstreamStatusRepository returns UstreamStatus repository
func (c *RepositoryContainer) UstreamStatusRepository() entity.UstreamStatusRepository {
	return entity.NewUstreamStatusRepository(dao.NewDatastoreHandler(c.ctx))
}

// LineNotificationRepository returns LineNotification repository
func (c *RepositoryContainer) LineNotificationRepository() entity.LineNotificationRepository {
	return entity.NewLineNotificationRepository(dao.NewDatastoreHandler(c.ctx))
}
