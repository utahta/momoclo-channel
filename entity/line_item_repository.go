package entity

import "github.com/utahta/momoclo-channel/dao"

type (
	// LineItemRepository interface
	LineItemRepository interface {
		Exists(string) bool
		Find(string) (*LineItem, error)
		Save(*LineItem) error
		Tx(dao.PersistenceHandler) LineItemRepository
	}

	// LineItemRepository operates LineItem entity
	lineItemRepository struct {
		dao.PersistenceHandler
	}
)

// NewLineItemRepository returns the LineItemRepository
func NewLineItemRepository(h dao.PersistenceHandler) LineItemRepository {
	return &lineItemRepository{h}
}

// Exists exists line item
func (repo *lineItemRepository) Exists(id string) bool {
	_, err := repo.Find(id)
	return err == nil
}

// Find finds line item given id
func (repo *lineItemRepository) Find(id string) (*LineItem, error) {
	item := &LineItem{ID: id}
	return item, repo.Get(item)
}

// Save saves line item
func (repo *lineItemRepository) Save(item *LineItem) error {
	return repo.Put(item)
}

// Tx can be used in RunInTransaction
func (repo *lineItemRepository) Tx(h dao.PersistenceHandler) LineItemRepository {
	return NewLineItemRepository(h)
}
