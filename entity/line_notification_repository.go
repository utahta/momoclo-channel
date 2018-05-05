package entity

import (
	"context"

	"github.com/utahta/momoclo-channel/dao"
)

type (
	// LineNotificationRepository interface
	LineNotificationRepository interface {
		FindAll(context.Context) ([]*LineNotification, error)
		Save(context.Context, *LineNotification) error
		Delete(context.Context, string) error
	}

	// lineNotificationRepository operates LineNotification entity
	lineNotificationRepository struct {
		dao.PersistenceHandler
	}
)

// NewLineNotificationRepository returns the LineNotificationRepository
func NewLineNotificationRepository(h dao.PersistenceHandler) LineNotificationRepository {
	return &lineNotificationRepository{h}
}

// FindAll finds all line notification entities
func (repo *lineNotificationRepository) FindAll(ctx context.Context) ([]*LineNotification, error) {
	kind := repo.Kind(ctx, &LineNotification{})
	q := repo.NewQuery(kind)

	var dst []*LineNotification
	return dst, repo.GetAll(ctx, q, &dst)
}

// Save saves given line notification entity
func (repo *lineNotificationRepository) Save(ctx context.Context, item *LineNotification) error {
	return repo.Put(ctx, item)
}

// Delete deletes given line notification entity
func (repo *lineNotificationRepository) Delete(ctx context.Context, id string) error {
	return repo.PersistenceHandler.Delete(ctx, &LineNotification{ID: id})
}
