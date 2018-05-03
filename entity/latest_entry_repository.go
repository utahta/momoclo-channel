package entity

import (
	"github.com/utahta/momoclo-channel/dao"
)

type (
	// LatestEntryRepository interface
	LatestEntryRepository interface {
		Save(*LatestEntry) error
		FindOrNewByURL(string, string) (*LatestEntry, error)
		GetURL(string) string
	}

	latestEntryRepository struct {
		dao.PersistenceHandler
	}
)

// NewLatestEntryRepository returns the LatestEntryRepository
func NewLatestEntryRepository(h dao.PersistenceHandler) LatestEntryRepository {
	return &latestEntryRepository{h}
}

// Save saves LatestEntry
func (repo *latestEntryRepository) Save(l *LatestEntry) error {
	return repo.Put(l)
}

// FindOrNewByURL finds LatestEntry given url
// if not found, returns new LatestEntry
func (repo *latestEntryRepository) FindOrNewByURL(code, urlStr string) (*LatestEntry, error) {
	l, err := NewLatestEntry(code, urlStr)
	if err != nil {
		return nil, err
	}

	err = repo.Get(l)
	if err == dao.ErrNoSuchEntity {
		return l, nil
	}
	if err != nil {
		return nil, err
	}
	return l, nil
}

// GetURL returns URL given code
func (repo *latestEntryRepository) GetURL(code string) string {
	l := &LatestEntry{ID: code}
	if err := repo.Get(l); err != nil {
		return ""
	}
	return l.URL
}
