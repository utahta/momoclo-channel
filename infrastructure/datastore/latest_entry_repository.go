package datastore

import (
	"context"

	"github.com/mjibson/goon"
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/domain"
	"github.com/utahta/momoclo-channel/domain/entity"
	"google.golang.org/appengine/datastore"
)

type LatestEntryRepository struct{}

func NewLatestEntryRepository() *LatestEntryRepository {
	return &LatestEntryRepository{}
}

func (repo *LatestEntryRepository) Save(c context.Context, l *entity.LatestEntry) error {
	g := goon.FromContext(c)
	return g.RunInTransaction(func(g *goon.Goon) error {
		_, err := g.Put(l)
		return err
	}, nil)
}

func (repo *LatestEntryRepository) FindByURL(c context.Context, urlStr string) (*entity.LatestEntry, error) {
	const errTag = "LatestEntryRepository.FindByURL failed"

	l := &entity.LatestEntry{ID: entity.ParseLatestEntryCode(urlStr)}
	g := goon.FromContext(c)

	err := g.Get(l)
	if err == datastore.ErrNoSuchEntity {
		return nil, domain.ErrNoSuchEntity
	}
	if err != nil {
		return nil, errors.Wrap(err, errTag)
	}
	return l, nil
}

func (repo *LatestEntryRepository) getURL(ctx context.Context, code string) string {
	g := goon.FromContext(ctx)
	l := &entity.LatestEntry{ID: code}
	if err := g.Get(l); err != nil {
		return ""
	}
	return l.URL
}

func (repo *LatestEntryRepository) GetTamaiURL(ctx context.Context) string {
	return repo.getURL(ctx, entity.LatestEntryCodeTamai)
}

func (repo *LatestEntryRepository) GetMomotaURL(ctx context.Context) string {
	return repo.getURL(ctx, entity.LatestEntryCodeMomota)
}

func (repo *LatestEntryRepository) GetAriyasuURL(ctx context.Context) string {
	return repo.getURL(ctx, entity.LatestEntryCodeAriyasu)
}

func (repo *LatestEntryRepository) GetSasakiURL(ctx context.Context) string {
	return repo.getURL(ctx, entity.LatestEntryCodeSasaki)
}

func (repo *LatestEntryRepository) GetTakagiURL(ctx context.Context) string {
	return repo.getURL(ctx, entity.LatestEntryCodeTakagi)
}

func (repo *LatestEntryRepository) GetHappycloURL(ctx context.Context) string {
	return repo.getURL(ctx, entity.LatestEntryCodeHappyclo)
}
