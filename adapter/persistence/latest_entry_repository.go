package persistence

import (
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/domain"
	"github.com/utahta/momoclo-channel/domain/entity"
	"github.com/utahta/momoclo-channel/domain/service/latest_entry"
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
func (repo *LatestEntryRepository) Save(l *entity.LatestEntry) error {
	return repo.Put(l)
}

// FindByURL finds LatestEntry given url
func (repo *LatestEntryRepository) FindByURL(urlStr string) (*entity.LatestEntry, error) {
	const errTag = "LatestEntryRepository.FindByURL failed"

	l, err := latest_entry.Parse(urlStr)
	if err != nil {
		return nil, errors.Wrap(err, errTag)
	}

	err = repo.Get(l)
	if err == domain.ErrNoSuchEntity {
		return nil, err
	}
	if err != nil {
		return nil, errors.Wrap(err, errTag)
	}
	return l, nil
}

func (repo *LatestEntryRepository) getURL(code string) string {
	l := &entity.LatestEntry{ID: code}
	if err := repo.Get(l); err != nil {
		return ""
	}
	return l.URL
}

func (repo *LatestEntryRepository) GetTamaiURL() string {
	return repo.getURL(entity.LatestEntryCodeTamai)
}

func (repo *LatestEntryRepository) GetMomotaURL() string {
	return repo.getURL(entity.LatestEntryCodeMomota)
}

func (repo *LatestEntryRepository) GetAriyasuURL() string {
	return repo.getURL(entity.LatestEntryCodeAriyasu)
}

func (repo *LatestEntryRepository) GetSasakiURL() string {
	return repo.getURL(entity.LatestEntryCodeSasaki)
}

func (repo *LatestEntryRepository) GetTakagiURL() string {
	return repo.getURL(entity.LatestEntryCodeTakagi)
}

func (repo *LatestEntryRepository) GetHappycloURL() string {
	return repo.getURL(entity.LatestEntryCodeHappyclo)
}
