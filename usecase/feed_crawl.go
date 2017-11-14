package usecase

import (
	"context"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/domain/core"
	"github.com/utahta/momoclo-channel/domain/event"
	"github.com/utahta/momoclo-channel/domain/model"
)

type (
	// FeedCrawl use case
	FeedCrawl struct {
		ctx             context.Context
		log             core.Logger
		feed            model.FeedFetcher
		taskQueue       event.TaskQueue
		latestEntryRepo model.LatestEntryRepository
	}

	// FeedCrawlParams input parameters
	FeedCrawlParams struct {
		Code string // target identify code
	}
)

// NewFeedCrawl returns Crawl use case
func NewFeedCrawl(
	ctx context.Context,
	log core.Logger,
	feed model.FeedFetcher,
	taskQueue event.TaskQueue,
	latestEntryRepo model.LatestEntryRepository) *FeedCrawl {
	return &FeedCrawl{
		ctx:             ctx,
		log:             log,
		feed:            feed,
		taskQueue:       taskQueue,
		latestEntryRepo: latestEntryRepo,
	}
}

// Do crawls a site and invokes tweet and line event
func (c *FeedCrawl) Do(params FeedCrawlParams) error {
	const errTag = "FeedCrawl.Do failed"

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
			event.Task{QueueName: "queue-tweet", Path: "/queue/feed/tweets/enqueue", Object: item},
			event.Task{QueueName: "queue-line", Path: "/queue/line", Object: item},
		)
	}
	if err := c.taskQueue.PushMulti(tasks); err != nil {
		return errors.Wrap(err, errTag)
	}
	return nil
}
