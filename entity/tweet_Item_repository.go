package entity

import (
	"context"

	"github.com/utahta/momoclo-channel/dao"
)

type (
	// TweetItemRepository interface
	TweetItemRepository interface {
		Exists(context.Context, string) bool
		Find(context.Context, string) (*TweetItem, error)
		Save(context.Context, *TweetItem) error
	}

	tweetItemRepository struct {
		dao.PersistenceHandler
	}
)

// NewTweetItemRepository returns the TweetItemRepository
func NewTweetItemRepository(h dao.PersistenceHandler) TweetItemRepository {
	return &tweetItemRepository{h}
}

// Exists exists tweet item
func (repo *tweetItemRepository) Exists(ctx context.Context, id string) bool {
	_, err := repo.Find(ctx, id)
	return err == nil
}

// Find finds tweet item given id
func (repo *tweetItemRepository) Find(ctx context.Context, id string) (*TweetItem, error) {
	item := &TweetItem{ID: id}
	return item, repo.Get(ctx, item)
}

// Save saves tweet item
func (repo *tweetItemRepository) Save(ctx context.Context, item *TweetItem) error {
	return repo.Put(ctx, item)
}
