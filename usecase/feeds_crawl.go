package usecase

import (
	"context"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/domain/core"
	"github.com/utahta/momoclo-channel/domain/model"
	"golang.org/x/sync/errgroup"
)

type (
	// FeedsCrawl use case
	FeedsCrawl struct {
		ctx   context.Context
		log   core.Logger
		crawl *FeedCrawl
	}
)

// NewFeedsCrawl returns CrawlAll use case
func NewFeedsCrawl(ctx context.Context, logger core.Logger, crawl *FeedCrawl) *FeedsCrawl {
	return &FeedsCrawl{
		ctx:   ctx,
		log:   logger,
		crawl: crawl,
	}
}

// Do crawls all sites
func (c *FeedsCrawl) Do() error {
	const errTag = "CrawlAll.Do failed"

	codes := []string{
		model.LatestEntryCodeMomota,
		model.LatestEntryCodeAriyasu,
		model.LatestEntryCodeTamai,
		model.LatestEntryCodeSasaki,
		model.LatestEntryCodeTakagi,
		model.LatestEntryCodeHappyclo,
		model.LatestEntryCodeAeNews,
		model.LatestEntryCodeYoutube,
	}

	eg := &errgroup.Group{}
	for _, code := range codes {
		code := code

		eg.Go(func() error {
			return c.crawl.Do(FeedCrawlParams{code})
		})
	}

	if err := eg.Wait(); err != nil {
		return errors.Wrap(err, errTag)
	}
	return nil
}
