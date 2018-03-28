package persistence

import (
	"github.com/utahta/momoclo-channel/types"
)

// TweetItemRepository operates TweetItem entity
type TweetItemRepository struct {
	types.PersistenceHandler
}

// NewTweetItemRepository returns the TweetItemRepository
func NewTweetItemRepository(h types.PersistenceHandler) types.TweetItemRepository {
	return &TweetItemRepository{h}
}

// Exists exists tweet item
func (repo *TweetItemRepository) Exists(id string) bool {
	_, err := repo.Find(id)
	return err == nil
}

// Find finds tweet item given id
func (repo *TweetItemRepository) Find(id string) (*types.TweetItem, error) {
	item := &types.TweetItem{ID: id}
	return item, repo.Get(item)
}

// Save saves tweet item
func (repo *TweetItemRepository) Save(item *types.TweetItem) error {
	return repo.Put(item)
}

// Tx can be used in RunInTransaction
func (repo *TweetItemRepository) Tx(h types.PersistenceHandler) types.TweetItemRepository {
	return NewTweetItemRepository(h)
}
