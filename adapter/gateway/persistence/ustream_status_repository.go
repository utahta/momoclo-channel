package persistence

import "github.com/utahta/momoclo-channel/types"

// UstreamStatusRepository operates UstreamStatus entity
type UstreamStatusRepository struct {
	types.PersistenceHandler
}

// NewUstreamStatusRepository returns the UstreamStatusRepository
func NewUstreamStatusRepository(h types.PersistenceHandler) types.UstreamStatusRepository {
	return &UstreamStatusRepository{h}
}

// Find finds ustream status entity
func (repo *UstreamStatusRepository) Find(id string) (*types.UstreamStatus, error) {
	entity := types.NewUstreamStatus()
	return entity, repo.Get(entity)
}

// Save saves ustream status entity
func (repo *UstreamStatusRepository) Save(item *types.UstreamStatus) error {
	return repo.Put(item)
}
