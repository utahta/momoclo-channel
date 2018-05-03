package entity

import "github.com/utahta/momoclo-channel/dao"

type (
	// LineNotificationRepository interface
	LineNotificationRepository interface {
		FindAll() ([]*LineNotification, error)
		Save(*LineNotification) error
		Delete(string) error
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
func (repo *lineNotificationRepository) FindAll() ([]*LineNotification, error) {
	kind := repo.Kind(&LineNotification{})
	q := repo.NewQuery(kind)

	var dst []*LineNotification
	return dst, repo.GetAll(q, &dst)
}

// Save saves given line notification entity
func (repo *lineNotificationRepository) Save(item *LineNotification) error {
	return repo.Put(item)
}

// Delete deletes given line notification entity
func (repo *lineNotificationRepository) Delete(id string) error {
	return repo.PersistenceHandler.Delete(&LineNotification{ID: id})
}
