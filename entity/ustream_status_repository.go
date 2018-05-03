package entity

import "github.com/utahta/momoclo-channel/types"

type (
	// UstreamStatusRepository operates UstreamStatus entity
	UstreamStatusRepository interface {
		Find(string) (*UstreamStatus, error)
		Save(*UstreamStatus) error
	}

	ustreamStatusRepository struct {
		types.PersistenceHandler
	}
)

// NewUstreamStatusRepository returns the UstreamStatusRepository
func NewUstreamStatusRepository(h types.PersistenceHandler) UstreamStatusRepository {
	return &ustreamStatusRepository{h}
}

// Find finds ustream status entity
func (repo *ustreamStatusRepository) Find(id string) (*UstreamStatus, error) {
	entity := NewUstreamStatus()
	return entity, repo.Get(entity)
}

// Save saves ustream status entity
func (repo *ustreamStatusRepository) Save(item *UstreamStatus) error {
	return repo.Put(item)
}
