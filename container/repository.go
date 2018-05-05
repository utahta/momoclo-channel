package container

import (
	"github.com/utahta/momoclo-channel/dao"
	"github.com/utahta/momoclo-channel/entity"
)

// RepositoryContainer dependency injection
type RepositoryContainer struct {
}

// Repository returns container of repositories
func Repository() *RepositoryContainer {
	return &RepositoryContainer{}
}

// LatestEntryRepository returns LatestEntry repository
func (c *RepositoryContainer) LatestEntryRepository() entity.LatestEntryRepository {
	return entity.NewLatestEntryRepository(dao.NewDatastoreHandler())
}

// TweetItemRepository returns TweetItem repository
func (c *RepositoryContainer) TweetItemRepository() entity.TweetItemRepository {
	return entity.NewTweetItemRepository(dao.NewDatastoreHandler())
}

// LineItemRepository returns LineItem repository
func (c *RepositoryContainer) LineItemRepository() entity.LineItemRepository {
	return entity.NewLineItemRepository(dao.NewDatastoreHandler())
}

// ReminderRepository returns Reminder repository
func (c *RepositoryContainer) ReminderRepository() entity.ReminderRepository {
	return entity.NewReminderRepository(dao.NewDatastoreHandler())
}

// UstreamStatusRepository returns UstreamStatus repository
func (c *RepositoryContainer) UstreamStatusRepository() entity.UstreamStatusRepository {
	return entity.NewUstreamStatusRepository(dao.NewDatastoreHandler())
}

// LineNotificationRepository returns LineNotification repository
func (c *RepositoryContainer) LineNotificationRepository() entity.LineNotificationRepository {
	return entity.NewLineNotificationRepository(dao.NewDatastoreHandler())
}
