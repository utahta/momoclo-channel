package entity

import "github.com/utahta/momoclo-channel/dao"

type (
	// TweetItemRepository interface
	TweetItemRepository interface {
		Exists(string) bool
		Find(string) (*TweetItem, error)
		Save(*TweetItem) error
		Tx(dao.PersistenceHandler) TweetItemRepository
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
func (repo *tweetItemRepository) Exists(id string) bool {
	_, err := repo.Find(id)
	return err == nil
}

// Find finds tweet item given id
func (repo *tweetItemRepository) Find(id string) (*TweetItem, error) {
	item := &TweetItem{ID: id}
	return item, repo.Get(item)
}

// Save saves tweet item
func (repo *tweetItemRepository) Save(item *TweetItem) error {
	return repo.Put(item)
}

// Tx can be used in RunInTransaction
func (repo *tweetItemRepository) Tx(h dao.PersistenceHandler) TweetItemRepository {
	return NewTweetItemRepository(h)
}
