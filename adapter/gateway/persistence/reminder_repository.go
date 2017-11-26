package persistence

import (
	"github.com/utahta/momoclo-channel/domain/model"
)

// ReminderRepository operates Reminder entity
type ReminderRepository struct {
	model.PersistenceHandler
}

// NewReminderRepository returns the ReminderRepository
func NewReminderRepository(h model.PersistenceHandler) model.ReminderRepository {
	return &ReminderRepository{h}
}

// FindAll finds all reminder entities
func (repo *ReminderRepository) FindAll() ([]*model.Reminder, error) {
	kind := repo.Kind(&model.Reminder{})
	q := repo.NewQuery(kind).Filter("Enabled =", true)

	var dst []*model.Reminder
	return dst, repo.GetAll(q, &dst)
}

// Save saves reminder entity
func (repo *ReminderRepository) Save(item *model.Reminder) error {
	return repo.Put(item)
}
