package usecase

import (
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/domain/core"
	"github.com/utahta/momoclo-channel/domain/event"
	"github.com/utahta/momoclo-channel/domain/model"
)

type (
	// CrawlFeed use case
	CrawlFeed struct {
		log             core.Logger
		feed            model.FeedFetcher
		taskQueue       event.TaskQueue
		latestEntryRepo model.LatestEntryRepository
	}

	// CrawlFeedParams input parameters
	CrawlFeedParams struct {
		Code string // target identify code
	}
)

// NewCrawlFeed returns Crawl use case
func NewCrawlFeed(
	log core.Logger,
	feed model.FeedFetcher,
	taskQueue event.TaskQueue,
	latestEntryRepo model.LatestEntryRepository) *CrawlFeed {
	return &CrawlFeed{
		log:             log,
		feed:            feed,
		taskQueue:       taskQueue,
		latestEntryRepo: latestEntryRepo,
	}
}

// Do crawls a site and invokes tweet and line event
func (c *CrawlFeed) Do(params CrawlFeedParams) error {
	const errTag = "CrawlFeed.Do failed"

	items, err := c.feed.Fetch(params.Code, 1, c.latestEntryRepo.GetURL(params.Code))
	if err != nil {
		return errors.Wrap(err, errTag)
	}
	if len(items) == 0 {
		return nil
	}

	// update latest entry
	item := items[0] // first item is the latest entry
	l, err := c.latestEntryRepo.FindOrCreateByURL(item.EntryURL)
	if err != nil {
		return errors.Wrapf(err, "%v: url:%v", errTag, item.EntryURL)
	}
	if !l.CreatedAt.IsZero() && l.URL == item.EntryURL {
		return nil // already get feeds. nothing to do
	}
	l.URL = item.EntryURL
	if err := c.latestEntryRepo.Save(l); err != nil {
		return errors.Wrapf(err, errTag)
	}

	// push events
	var tasks []event.Task
	for _, item := range items {
		tasks = append(tasks,
			event.Task{QueueName: "enqueue", Path: "/enqueue/tweets", Object: item},
			event.Task{QueueName: "enqueue", Path: "/enqueue/lines", Object: item},
		)
	}
	if err := c.taskQueue.PushMulti(tasks); err != nil {
		return errors.Wrap(err, errTag)
	}
	c.log.Infof("crawl feed items:%v", items)

	return nil
}
