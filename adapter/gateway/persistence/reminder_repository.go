package persistence

import (
	"github.com/utahta/momoclo-channel/types"
)

// ReminderRepository operates Reminder entity
type ReminderRepository struct {
	types.PersistenceHandler
}

// NewReminderRepository returns the ReminderRepository
func NewReminderRepository(h types.PersistenceHandler) types.ReminderRepository {
	return &ReminderRepository{h}
}

// FindAll finds all reminder entities
func (repo *ReminderRepository) FindAll() ([]*types.Reminder, error) {
	kind := repo.Kind(&types.Reminder{})
	q := repo.NewQuery(kind).Filter("Enabled =", true)

	var dst []*types.Reminder
	return dst, repo.GetAll(q, &dst)
}

// Save saves reminder entity
func (repo *ReminderRepository) Save(item *types.Reminder) error {
	return repo.Put(item)
}
