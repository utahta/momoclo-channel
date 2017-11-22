package persistence

import (
	"github.com/utahta/momoclo-channel/domain/model"
)

// LineItemRepository operates LineItem entity
type LineItemRepository struct {
	model.PersistenceHandler
}

// NewLineItemRepository returns the LineItemRepository
func NewLineItemRepository(h model.PersistenceHandler) *LineItemRepository {
	return &LineItemRepository{h}
}

// Exists exists line item
func (repo *LineItemRepository) Exists(id string) bool {
	_, err := repo.Find(id)
	return err == nil
}

// Find finds line item given id
func (repo *LineItemRepository) Find(id string) (*model.LineItem, error) {
	item := &model.LineItem{ID: id}
	return item, repo.Get(item)
}

// Save saves line item
func (repo *LineItemRepository) Save(item *model.LineItem) error {
	return repo.Put(item)
}
