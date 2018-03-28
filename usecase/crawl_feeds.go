package usecase

import (
	"time"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/log"
	"github.com/utahta/momoclo-channel/timeutil"
	"github.com/utahta/momoclo-channel/types"
	"golang.org/x/sync/errgroup"
)

type (
	// CrawlFeeds use case
	CrawlFeeds struct {
		log   log.Logger
		crawl *CrawlFeed
	}
)

// NewCrawlFeeds returns CrawlAll use case
func NewCrawlFeeds(logger log.Logger, crawl *CrawlFeed) *CrawlFeeds {
	return &CrawlFeeds{
		log:   logger,
		crawl: crawl,
	}
}

// Do crawls all known sites
func (c *CrawlFeeds) Do() error {
	const errTag = "CrawlFeeds.Do failed"

	now := timeutil.Now()
	codes := []types.FeedCode{
		types.FeedCodeMomota,
		types.FeedCodeTamai,
		types.FeedCodeSasaki,
		types.FeedCodeTakagi,
		types.FeedCodeAeNews,
		types.FeedCodeYoutube,
	}
	if now.Weekday() == time.Sunday {
		codes = append(codes, types.FeedCodeHappyclo)
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
