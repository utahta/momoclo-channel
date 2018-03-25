package usecase

import (
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/domain/core"
	"github.com/utahta/momoclo-channel/domain/model"
	"github.com/utahta/momoclo-channel/event"
	"github.com/utahta/momoclo-channel/event/eventtask"
	"github.com/utahta/momoclo-channel/log"
)

type (
	// CrawlFeed use case
	CrawlFeed struct {
		log       log.Logger
		feed      model.FeedFetcher
		taskQueue event.TaskQueue
		repo      model.LatestEntryRepository
	}

	// CrawlFeedParams input parameters
	CrawlFeedParams struct {
		Code model.FeedCode // target identify code
	}
)

// NewCrawlFeed returns Crawl use case
func NewCrawlFeed(
	log log.Logger,
	feed model.FeedFetcher,
	taskQueue event.TaskQueue,
	repo model.LatestEntryRepository) *CrawlFeed {
	return &CrawlFeed{
		log:       log,
		feed:      feed,
		taskQueue: taskQueue,
		repo:      repo,
	}
}

// Do crawls a site and invokes tweet and line event
func (use *CrawlFeed) Do(params CrawlFeedParams) error {
	const errTag = "CrawlFeed.Do failed"

	items, err := use.feed.Fetch(params.Code, 1, use.repo.GetURL(params.Code.String()))
	if err != nil {
		return errors.Wrap(err, errTag)
	}
	if len(items) == 0 {
		return nil
	}
	for i := range items {
		if err := core.Validate(items[i]); err != nil {
			use.log.Errorf("%v: validate error i:%v items:%v err:%v", errTag, i, items, err)
			return errors.Wrap(err, errTag)
		}
	}

	// update latest entry
	item := items[0] // first item is the latest entry
	l, err := use.repo.FindOrNewByURL(item.EntryURL)
	if err != nil {
		return errors.Wrapf(err, "%v: url:%v", errTag, item.EntryURL)
	}
	if l.URL == item.EntryURL && l.PublishedAt.Equal(item.PublishedAt) {
		return nil // already get feeds. nothing to do
	}
	l.URL = item.EntryURL
	l.PublishedAt = item.PublishedAt
	if err := use.repo.Save(l); err != nil {
		return errors.Wrapf(err, errTag)
	}

	// push events
	var tasks []event.Task
	for _, item := range items {
		tasks = append(tasks,
			eventtask.NewEnqueueTweets(item),
			eventtask.NewEnqueueLines(item),
		)
	}
	if err := use.taskQueue.PushMulti(tasks); err != nil {
		return errors.Wrap(err, errTag)
	}
	use.log.Infof("crawl feed items:%v", items)

	return nil
}
