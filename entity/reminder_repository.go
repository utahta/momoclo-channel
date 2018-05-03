package entity

import "github.com/utahta/momoclo-channel/dao"

type (
	// ReminderRepository interface
	ReminderRepository interface {
		FindAll() ([]*Reminder, error)
		Save(*Reminder) error
	}

	reminderRepository struct {
		dao.PersistenceHandler
	}
)

// NewReminderRepository returns the ReminderRepository
func NewReminderRepository(h dao.PersistenceHandler) ReminderRepository {
	return &reminderRepository{h}
}

// FindAll finds all reminder entities
func (repo *reminderRepository) FindAll() ([]*Reminder, error) {
	kind := repo.Kind(&Reminder{})
	q := repo.NewQuery(kind).Filter("Enabled =", true)

	var dst []*Reminder
	return dst, repo.GetAll(q, &dst)
}

// Save saves reminder entity
func (repo *reminderRepository) Save(item *Reminder) error {
	return repo.Put(item)
}
