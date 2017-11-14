package usecase

import (
	"context"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/domain/core"
	"github.com/utahta/momoclo-channel/domain/entity"
	"golang.org/x/sync/errgroup"
)

type (
	// CrawlAll use case
	CrawlAll struct {
		ctx   context.Context
		log   core.Logger
		crawl *Crawl
	}
)

// NewCrawlAll returns CrawlAll use case
func NewCrawlAll(ctx context.Context, logger core.Logger, crawl *Crawl) *CrawlAll {
	return &CrawlAll{
		ctx:   ctx,
		log:   logger,
		crawl: crawl,
	}
}

// Do crawls all sites
func (c *CrawlAll) Do() error {
	const errTag = "CrawlAll.Do failed"

	codes := []string{
		entity.LatestEntryCodeMomota,
		entity.LatestEntryCodeAriyasu,
		entity.LatestEntryCodeTamai,
		entity.LatestEntryCodeSasaki,
		entity.LatestEntryCodeTakagi,
		entity.LatestEntryCodeHappyclo,
		entity.LatestEntryCodeAeNews,
		entity.LatestEntryCodeYoutube,
	}

	eg := &errgroup.Group{}
	for _, code := range codes {
		code := code

		eg.Go(func() error {
			return c.crawl.Do(CrawlParams{code})
		})
	}

	if err := eg.Wait(); err != nil {
		return errors.Wrap(err, errTag)
	}
	return nil
}
