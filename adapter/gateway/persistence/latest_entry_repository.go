package persistence

import (
	"github.com/utahta/momoclo-channel/types"
)

// LatestEntryRepository operates LatestEntry entity
type LatestEntryRepository struct {
	types.PersistenceHandler
}

// NewLatestEntryRepository returns the LatestEntryRepository
func NewLatestEntryRepository(h types.PersistenceHandler) types.LatestEntryRepository {
	return &LatestEntryRepository{h}
}

// Save saves LatestEntry
func (repo *LatestEntryRepository) Save(l *types.LatestEntry) error {
	return repo.Put(l)
}

// FindOrNewByURL finds LatestEntry given url
// if not found, returns new LatestEntry
func (repo *LatestEntryRepository) FindOrNewByURL(urlStr string) (*types.LatestEntry, error) {
	l, err := types.NewLatestEntry(urlStr)
	if err != nil {
		return nil, err
	}

	err = repo.Get(l)
	if err == types.ErrNoSuchEntity {
		return l, nil
	}
	if err != nil {
		return nil, err
	}
	return l, nil
}

// GetURL returns URL given code
func (repo *LatestEntryRepository) GetURL(code string) string {
	l := &types.LatestEntry{ID: code}
	if err := repo.Get(l); err != nil {
		return ""
	}
	return l.URL
}
