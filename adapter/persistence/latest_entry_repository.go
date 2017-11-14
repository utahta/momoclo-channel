package persistence

import (
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/domain"
	"github.com/utahta/momoclo-channel/domain/entity"
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
func (repo *LatestEntryRepository) Save(l *entity.LatestEntry) error {
	return repo.Put(l)
}

// FindOrCreateByURL finds LatestEntry given url
// if not found returns new LatestEntry
func (repo *LatestEntryRepository) FindOrCreateByURL(urlStr string) (*entity.LatestEntry, error) {
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
	l := &entity.LatestEntry{ID: code}
	if err := repo.Get(l); err != nil {
		return ""
	}
	return l.URL
}

// GetTamaiURL returns Shiori Tamai blog url
func (repo *LatestEntryRepository) GetTamaiURL() string {
	return repo.GetURL(entity.LatestEntryCodeTamai)
}

// GetMomotaURL returns Kanako Momota blog url
func (repo *LatestEntryRepository) GetMomotaURL() string {
	return repo.GetURL(entity.LatestEntryCodeMomota)
}

// GetAriyasuURL returns Momoka Ariyasu blog url
func (repo *LatestEntryRepository) GetAriyasuURL() string {
	return repo.GetURL(entity.LatestEntryCodeAriyasu)
}

// GetSasakiURL returns Ayaka Sasaki blog url
func (repo *LatestEntryRepository) GetSasakiURL() string {
	return repo.GetURL(entity.LatestEntryCodeSasaki)
}

// GetTakagiURL returns Reni Takagi blog url
func (repo *LatestEntryRepository) GetTakagiURL() string {
	return repo.GetURL(entity.LatestEntryCodeTakagi)
}

// GetHappycloURL returns happyclo site url
func (repo *LatestEntryRepository) GetHappycloURL() string {
	return repo.GetURL(entity.LatestEntryCodeHappyclo)
}
