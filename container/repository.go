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

// TweetItemRepository returns TweetItem repository
func (c *RepositoryContainer) TweetItemRepository() model.TweetItemRepository {
	return persistence.NewTweetItemRepository(dao.NewDatastoreHandler(c.ctx))
}

// ReminderRepository returns Reminder repository
func (c *RepositoryContainer) ReminderRepository() model.ReminderRepository {
	return persistence.NewReminderRepository(dao.NewDatastoreHandler(c.ctx))
}

// UstreamStatusRepository returns UstreamStatus repository
func (c *RepositoryContainer) UstreamStatusRepository() model.UstreamStatusRepository {
	return persistence.NewUstreamStatusRepository(dao.NewDatastoreHandler(c.ctx))
}

// LineNotificationRepository returns LineNotification repository
func (c *RepositoryContainer) LineNotificationRepository() model.LineNotificationRepository {
	return persistence.NewLineNotificationRepository(dao.NewDatastoreHandler(c.ctx))
}
