package usecase

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/domain/core"
	"github.com/utahta/momoclo-channel/domain/entity"
	"github.com/utahta/momoclo-channel/domain/event"
)

type (
	// Crawler interface
	Crawler interface {
		Fetch(code string, maxItemNum int, latestURL string) ([]*CrawlItem, error)
	}

	// CrawlItem crawled item
	CrawlItem struct {
		Title       string
		URL         string
		EntryTitle  string
		EntryURL    string
		ImageURLs   []string
		VideoURLs   []string
		PublishedAt time.Time
	}

	// Crawl use case
	Crawl struct {
		ctx             context.Context
		log             core.Logger
		crawler         Crawler
		taskQueue       event.TaskQueue
		latestEntryRepo entity.LatestEntryRepository
	}

	// CrawlParams input parameters
	CrawlParams struct {
		Code string // target identify code
	}
)

// NewCrawl returns Crawl use case
func NewCrawl(
	ctx context.Context,
	log core.Logger,
	crawler Crawler,
	taskQueue event.TaskQueue,
	latestEntryRepo entity.LatestEntryRepository) *Crawl {
	return &Crawl{
		ctx:             ctx,
		log:             log,
		crawler:         crawler,
		taskQueue:       taskQueue,
		latestEntryRepo: latestEntryRepo,
	}
}

// Do crawls a web site and invokes tweet and line event
func (c *Crawl) Do(params CrawlParams) error {
	const errTag = "Crawl.Do failed"

	items, err := c.crawler.Fetch(params.Code, 1, c.latestEntryRepo.GetURL(params.Code))
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
		c.log.Warningf("%v: url:%v err:%v", errTag, item.EntryURL, err)
	} else {
		if l.CreatedAt.IsZero() || l.URL != item.EntryURL {
			if err := c.latestEntryRepo.Save(l); err != nil {
				c.log.Warningf("%v: err:%v", errTag, err)
			}
		}
	}

	// push events
	var tasks []event.Task
	for _, item := range items {
		tasks = append(tasks,
			event.Task{QueueName: "queue-tweet", Path: "/queue/tweet", Object: item},
			event.Task{QueueName: "queue-line", Path: "/queue/line", Object: item},
		)
	}
	if err := c.taskQueue.PushMulti(tasks); err != nil {
		return errors.Wrap(err, errTag)
	}
	return nil
}
