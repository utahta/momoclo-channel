package entity

import (
	"context"

	"github.com/utahta/momoclo-channel/dao"
)

type (
	// LatestEntryRepository interface
	LatestEntryRepository interface {
		Save(context.Context, *LatestEntry) error
		FindOrNewByURL(context.Context, string, string) (*LatestEntry, error)
		GetURL(context.Context, string) string
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
func (repo *latestEntryRepository) Save(ctx context.Context, l *LatestEntry) error {
	return repo.Put(ctx, l)
}

// FindOrNewByURL finds LatestEntry given url
// if not found, returns new LatestEntry
func (repo *latestEntryRepository) FindOrNewByURL(ctx context.Context, code, urlStr string) (*LatestEntry, error) {
	l, err := NewLatestEntry(code, urlStr)
	if err != nil {
		return nil, err
	}

	err = repo.Get(ctx, l)
	if err == dao.ErrNoSuchEntity {
		return l, nil
	}
	if err != nil {
		return nil, err
	}
	return l, nil
}

// GetURL returns URL given code
func (repo *latestEntryRepository) GetURL(ctx context.Context, code string) string {
	l := &LatestEntry{ID: code}
	if err := repo.Get(ctx, l); err != nil {
		return ""
	}
	return l.URL
}
