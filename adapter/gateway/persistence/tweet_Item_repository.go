package persistence

import (
	"github.com/utahta/momoclo-channel/domain/model"
)

// TweetItemRepository operates TweetItem entity
type TweetItemRepository struct {
	model.PersistenceHandler
}

// NewTweetItemRepository returns the TweetItemRepository
func NewTweetItemRepository(h model.PersistenceHandler) model.TweetItemRepository {
	return &TweetItemRepository{h}
}

// Exists exists tweet item
func (repo *TweetItemRepository) Exists(id string) bool {
	_, err := repo.Find(id)
	return err == nil
}

// Find finds tweet item given id
func (repo *TweetItemRepository) Find(id string) (*model.TweetItem, error) {
	item := &model.TweetItem{ID: id}
	return item, repo.Get(item)
}

// Save saves tweet item
func (repo *TweetItemRepository) Save(item *model.TweetItem) error {
	return repo.Put(item)
}

// Tx can be used in RunInTransaction
func (repo *TweetItemRepository) Tx(h model.PersistenceHandler) model.TweetItemRepository {
	return NewTweetItemRepository(h)
}
