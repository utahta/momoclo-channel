package entity

import (
	"context"

	"github.com/utahta/momoclo-channel/dao"
)

type (
	// LineItemRepository interface
	LineItemRepository interface {
		Exists(context.Context, string) bool
		Find(context.Context, string) (*LineItem, error)
		Save(context.Context, *LineItem) error
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
func (repo *lineItemRepository) Exists(ctx context.Context, id string) bool {
	_, err := repo.Find(ctx, id)
	return err == nil
}

// Find finds line item given id
func (repo *lineItemRepository) Find(ctx context.Context, id string) (*LineItem, error) {
	item := &LineItem{ID: id}
	return item, repo.Get(ctx, item)
}

// Save saves line item
func (repo *lineItemRepository) Save(ctx context.Context, item *LineItem) error {
	return repo.Put(ctx, item)
}
