package persistence

import (
	"github.com/utahta/momoclo-channel/domain/model"
)

// UstreamStatusRepository operates UstreamStatus entity
type UstreamStatusRepository struct {
	model.PersistenceHandler
}

// NewUstreamStatusRepository returns the UstreamStatusRepository
func NewUstreamStatusRepository(h model.PersistenceHandler) model.UstreamStatusRepository {
	return &UstreamStatusRepository{h}
}

// Find finds ustream status entity
func (repo *UstreamStatusRepository) Find(id string) (*model.UstreamStatus, error) {
	entity := model.NewUstreamStatus()
	return entity, repo.Get(entity)
}

// Save saves ustream status entity
func (repo *UstreamStatusRepository) Save(item *model.UstreamStatus) error {
	return repo.Put(item)
}
