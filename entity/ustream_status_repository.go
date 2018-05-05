package entity

import (
	"context"

	"github.com/utahta/momoclo-channel/dao"
)

type (
	// UstreamStatusRepository operates UstreamStatus entity
	UstreamStatusRepository interface {
		Find(context.Context, string) (*UstreamStatus, error)
		Save(context.Context, *UstreamStatus) error
	}

	ustreamStatusRepository struct {
		dao.PersistenceHandler
	}
)

// NewUstreamStatusRepository returns the UstreamStatusRepository
func NewUstreamStatusRepository(h dao.PersistenceHandler) UstreamStatusRepository {
	return &ustreamStatusRepository{h}
}

// Find finds ustream status entity
func (repo *ustreamStatusRepository) Find(ctx context.Context, id string) (*UstreamStatus, error) {
	entity := NewUstreamStatus()
	return entity, repo.Get(ctx, entity)
}

// Save saves ustream status entity
func (repo *ustreamStatusRepository) Save(ctx context.Context, item *UstreamStatus) error {
	return repo.Put(ctx, item)
}
