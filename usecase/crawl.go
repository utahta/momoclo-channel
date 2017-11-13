package usecase

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/domain"
	"github.com/utahta/momoclo-channel/domain/core"
	"github.com/utahta/momoclo-channel/domain/entity"
	"github.com/utahta/momoclo-channel/domain/event"
	"github.com/utahta/momoclo-channel/domain/service/latestentry"
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

	// Crawl crawls a web site use case
	Crawl struct {
		ctx             context.Context
		log             core.Logger
		crawler         Crawler
		event           event.TaskQueue
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
	event event.TaskQueue,
	latestEntryRepo entity.LatestEntryRepository) *Crawl {
	return &Crawl{
		ctx:             ctx,
		log:             log,
		crawler:         crawler,
		event:           event,
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

	c.updateLatestEntry(items)

	for _, item := range items {
		if err := c.event.Push(event.Task{QueueName: "queue-tweet", Path: "/queue/tweet", Object: item}); err != nil {
			c.log.Errorf("%v: queue tweet err:%v %#v", errTag, err, item)
			continue
		}

		if err := c.event.Push(event.Task{QueueName: "queue-line", Path: "/queue/line", Object: item}); err != nil {
			c.log.Errorf("%v: queue line err:%v %#v", errTag, err, item)
			continue
		}
	}
	return nil
}

func (c *Crawl) updateLatestEntry(items []*CrawlItem) {
	const errTag = "Crawl.updateLatestEntry failed"

	if len(items) == 0 {
		return
	}

	item := items[0] // first item is latest entry
	l, err := c.latestEntryRepo.FindByURL(item.EntryURL)
	if err != nil {
		if err == domain.ErrNoSuchEntity {
			l, err = latestentry.Parse(item.EntryURL)
			if err != nil {
				c.log.Warningf("%v: parse url:%v err:%v", errTag, item.EntryURL, err)
				return
			}
		} else {
			c.log.Errorf("%v: FindByURL url:%v err:%v", errTag, item.EntryURL, err)
			return
		}
	} else {
		if l.URL == item.EntryURL {
			return
		}
	}

	if err := c.latestEntryRepo.Save(l); err != nil {
		c.log.Warningf("%v: put latest entry. err:%v", errTag, err)
		return
	}
}
