package persistence

import (
	"github.com/utahta/momoclo-channel/types"
)

// LineItemRepository operates LineItem entity
type LineItemRepository struct {
	types.PersistenceHandler
}

// NewLineItemRepository returns the LineItemRepository
func NewLineItemRepository(h types.PersistenceHandler) *LineItemRepository {
	return &LineItemRepository{h}
}

// Exists exists line item
func (repo *LineItemRepository) Exists(id string) bool {
	_, err := repo.Find(id)
	return err == nil
}

// Find finds line item given id
func (repo *LineItemRepository) Find(id string) (*types.LineItem, error) {
	item := &types.LineItem{ID: id}
	return item, repo.Get(item)
}

// Save saves line item
func (repo *LineItemRepository) Save(item *types.LineItem) error {
	return repo.Put(item)
}

// Tx can be used in RunInTransaction
func (repo *LineItemRepository) Tx(h types.PersistenceHandler) types.LineItemRepository {
	return NewLineItemRepository(h)
}
