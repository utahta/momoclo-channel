package persistence

import (
	"github.com/utahta/momoclo-channel/domain/model"
)

// LineNotificationRepository operates LineNotification entity
type LineNotificationRepository struct {
	model.PersistenceHandler
}

// NewLineNotificationRepository returns the LineNotificationRepository
func NewLineNotificationRepository(h model.PersistenceHandler) model.LineNotificationRepository {
	return &LineNotificationRepository{h}
}

// FindAll finds all line notification entities
func (repo *LineNotificationRepository) FindAll() ([]*model.LineNotification, error) {
	kind := repo.Kind(&model.LineNotification{})
	q := repo.NewQuery(kind)

	var dst []*model.LineNotification
	return dst, repo.GetAll(q, &dst)
}

// Save saves given line notification entity
func (repo *LineNotificationRepository) Save(item *model.LineNotification) error {
	return repo.Put(item)
}

// Delete deletes given line notification entity
func (repo *LineNotificationRepository) Delete(id string) error {
	return repo.PersistenceHandler.Delete(&model.LineNotification{ID: id})
}
