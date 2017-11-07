package persistence

import (
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/domain"
	"github.com/utahta/momoclo-channel/domain/entity"
)

type LatestEntryRepository struct {
	DatastoreHandler
}

func NewLatestEntryRepository(h DatastoreHandler) *LatestEntryRepository {
	return &LatestEntryRepository{h}
}

func (repo *LatestEntryRepository) Save(l *entity.LatestEntry) error {
	return repo.RunInTransaction(func(h DatastoreHandler) error {
		return h.Put(l)
	}, nil)
}

// FindByURL finds LatestEntry given url
func (repo *LatestEntryRepository) FindByURL(urlStr string) (*entity.LatestEntry, error) {
	const errTag = "LatestEntryRepository.FindByURL failed"

	code, err := entity.ParseLatestEntryCode(urlStr)
	if err != nil {
		return nil, errors.Wrap(err, errTag)
	}

	l := &entity.LatestEntry{ID: code}
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
