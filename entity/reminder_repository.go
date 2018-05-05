package entity

import (
	"context"

	"github.com/utahta/momoclo-channel/dao"
)

type (
	// ReminderRepository interface
	ReminderRepository interface {
		FindAll(context.Context) ([]*Reminder, error)
		Save(context.Context, *Reminder) error
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
func (repo *reminderRepository) FindAll(ctx context.Context) ([]*Reminder, error) {
	kind := repo.Kind(ctx, &Reminder{})
	q := repo.NewQuery(kind).Filter("Enabled =", true)

	var dst []*Reminder
	return dst, repo.GetAll(ctx, q, &dst)
}

// Save saves reminder entity
func (repo *reminderRepository) Save(ctx context.Context, item *Reminder) error {
	return repo.Put(ctx, item)
}
