package usecase

import (
	"context"

	"github.com/utahta/momoclo-channel/domain"
	"github.com/utahta/momoclo-channel/domain/entity"
	"github.com/utahta/momoclo-channel/domain/event"
	"github.com/utahta/momoclo-channel/domain/service/latestentry"
	"github.com/utahta/momoclo-crawler"
	"google.golang.org/appengine/urlfetch"
)

type (
	// CrawlOne crawls a web site use case
	CrawlOne struct {
		ctx             context.Context
		log             event.Logger
		latestEntryRepo entity.LatestEntryRepository
	}

	// CrawlOneInputPort input parameters
	CrawlOneInputPort struct {
		Code string // target identify code
	}
)

// NewCrawlOne returns CrawlOne use case
func NewCrawlOne(ctx context.Context, log event.Logger, latestEntryRepo entity.LatestEntryRepository) *CrawlOne {
	return &CrawlOne{
		ctx:             ctx,
		log:             log,
		latestEntryRepo: latestEntryRepo,
	}
}

// Do crawls a web site and invokes tweet and line event
func (c *CrawlOne) Do(params CrawlOneInputPort) error {
	const errTag = "CrawlOne.Do failed"

	client := c.getCrawlerClient(params.Code)

	ch, err := client.Fetch()
	if err != nil {
		c.log.Error(err)
		return err
	}

	c.updateLatestEntry(ch)

	//if err := c.PushTweet(ctx, ch); err != nil {
	//	log.Errorf(ctx, "%v: push tweet queue. err:%v", errTag, err)
	//}
	//if err := c.PushLine(ctx, ch); err != nil {
	//	log.Errorf(ctx, "%v: push line queue. err:%v", errTag, err)
	//}
	return nil
}

func (c *CrawlOne) getCrawlerClient(code string) *crawler.ChannelClient {

	option := crawler.WithHTTPClient(urlfetch.Client(c.ctx))
	switch code {
	case entity.LatestEntryCodeTamai:
		cw, _ := crawler.NewTamaiBlogChannelClient(1, c.latestEntryRepo.GetTamaiURL(), option)
		return cw
	}

	return nil

	//option := crawler.WithHTTPClient(urlfetch.Client(ctx))
	//clients := []*crawler.ChannelClient{
	//	c.retrieveChannelClient(crawler.NewTamaiBlogChannelClient(1, c.latestEntryRepo.GetTamaiURL(), option)),
	//	c.retrieveChannelClient(crawler.NewMomotaBlogChannelClient(1, c.latestEntryRepo.GetMomotaURL(), option)),
	//	c.retrieveChannelClient(crawler.NewAriyasuBlogChannelClient(1, c.latestEntryRepo.GetAriyasuURL(), option)),
	//	c.retrieveChannelClient(crawler.NewSasakiBlogChannelClient(1, c.latestEntryRepo.GetSasakiURL(), option)),
	//	c.retrieveChannelClient(crawler.NewTakagiBlogChannelClient(1, c.latestEntryRepo.GetTakagiURL(), option)),
	//	c.retrieveChannelClient(crawler.NewAeNewsChannelClient(option)),
	//	c.retrieveChannelClient(crawler.NewYoutubeChannelClient(option)),
	//}

	//now := timeutil.Now().In(config.JST)
	//if (now.Weekday() == time.Sunday && now.Hour() == 16 && (now.Minute() >= 55 && now.Minute() <= 59)) ||
	//	(now.Hour() >= 8 && now.Hour() <= 23 && (now.Minute() == 0 || now.Minute() == 30)) {
	//	clients = append(clients, c.retrieveChannelClient(crawler.NewHappycloChannelClient(c.latestEntryRepo.GetHappycloURL(), option)))
	//}

	//return clients
}

func (c *CrawlOne) updateLatestEntry(ch *crawler.Channel) {
	const errTag = "CrawlOne.updateLatestEntry failed"

	for _, item := range ch.Items {
		l, err := c.latestEntryRepo.FindByURL(item.Url)
		if err == domain.ErrNoSuchEntity {
			l, err = latestentry.Parse(item.Url)
			if err != nil {
				c.log.Warningf("%v: parse url:%v err:%v", errTag, item.Url, err)
				continue
			}
		} else if err != nil {
			c.log.Errorf("%v: FindByURL url:%v err:%v", errTag, item.Url, err)
			continue
		} else {
			if l.URL == item.Url {
				continue
			}
		}

		if err := c.latestEntryRepo.Save(l); err != nil {
			c.log.Warningf("%v: put latest entry. err:%v", errTag, err)
			continue
		}
		break // first item equals latest item
	}
}
