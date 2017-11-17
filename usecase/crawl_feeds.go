package usecase

import (
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/domain/core"
	"github.com/utahta/momoclo-channel/domain/model"
	"golang.org/x/sync/errgroup"
)

type (
	// CrawlFeeds use case
	CrawlFeeds struct {
		log   core.Logger
		crawl *CrawlFeed
	}
)

// NewCrawlFeeds returns CrawlAll use case
func NewCrawlFeeds(logger core.Logger, crawl *CrawlFeed) *CrawlFeeds {
	return &CrawlFeeds{
		log:   logger,
		crawl: crawl,
	}
}

// Do crawls all known sites
func (c *CrawlFeeds) Do() error {
	const errTag = "CrawlFeeds.Do failed"

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
			return c.crawl.Do(CrawlFeedParams{code})
		})
	}

	if err := eg.Wait(); err != nil {
		return errors.Wrap(err, errTag)
	}
	return nil
}
