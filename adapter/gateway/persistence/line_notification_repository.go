package persistence

import (
	"github.com/utahta/momoclo-channel/types"
)

// LineNotificationRepository operates LineNotification entity
type LineNotificationRepository struct {
	types.PersistenceHandler
}

// NewLineNotificationRepository returns the LineNotificationRepository
func NewLineNotificationRepository(h types.PersistenceHandler) types.LineNotificationRepository {
	return &LineNotificationRepository{h}
}

// FindAll finds all line notification entities
func (repo *LineNotificationRepository) FindAll() ([]*types.LineNotification, error) {
	kind := repo.Kind(&types.LineNotification{})
	q := repo.NewQuery(kind)

	var dst []*types.LineNotification
	return dst, repo.GetAll(q, &dst)
}

// Save saves given line notification entity
func (repo *LineNotificationRepository) Save(item *types.LineNotification) error {
	return repo.Put(item)
}

// Delete deletes given line notification entity
func (repo *LineNotificationRepository) Delete(id string) error {
	return repo.PersistenceHandler.Delete(&types.LineNotification{ID: id})
}
