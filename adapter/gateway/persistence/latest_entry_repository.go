package persistence

import (
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/domain"
	"github.com/utahta/momoclo-channel/domain/model"
	"github.com/utahta/momoclo-channel/domain/service/latestentry"
)

// LatestEntryRepository operates datastore
type LatestEntryRepository struct {
	DatastoreHandler
}

// NewLatestEntryRepository returns the LatestEntryRepository
func NewLatestEntryRepository(h DatastoreHandler) *LatestEntryRepository {
	return &LatestEntryRepository{h}
}

// Save saves LatestEntry
func (repo *LatestEntryRepository) Save(l *model.LatestEntry) error {
	return repo.Put(l)
}

// FindOrCreateByURL finds LatestEntry given url
// if not found, returns new LatestEntry
func (repo *LatestEntryRepository) FindOrCreateByURL(urlStr string) (*model.LatestEntry, error) {
	const errTag = "LatestEntryRepository.FindOrCreateByURL failed"

	l, err := latestentry.Parse(urlStr)
	if err != nil {
		return nil, errors.Wrap(err, errTag)
	}

	err = repo.Get(l)
	if err == domain.ErrNoSuchEntity {
		return l, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, errTag)
	}
	return l, nil
}

// GetURL returns URL given code
func (repo *LatestEntryRepository) GetURL(code string) string {
	l := &model.LatestEntry{ID: code}
	if err := repo.Get(l); err != nil {
		return ""
	}
	return l.URL
}
