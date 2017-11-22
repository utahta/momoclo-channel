package container

import (
	"context"

	"github.com/utahta/momoclo-channel/adapter/gateway/api/customsearch"
	"github.com/utahta/momoclo-channel/adapter/gateway/api/linebot"
	"github.com/utahta/momoclo-channel/adapter/gateway/api/linenotify"
	"github.com/utahta/momoclo-channel/adapter/gateway/api/twitter"
	"github.com/utahta/momoclo-channel/adapter/gateway/api/ustream"
	"github.com/utahta/momoclo-channel/adapter/gateway/crawler"
	"github.com/utahta/momoclo-channel/infrastructure/dao"
	"github.com/utahta/momoclo-channel/infrastructure/event"
	"github.com/utahta/momoclo-channel/infrastructure/log"
	"github.com/utahta/momoclo-channel/usecase"
)

// UsecaseContainer dependency injection
type UsecaseContainer struct {
	ctx  context.Context
	repo *RepositoryContainer
}

// Usecase returns container of use case
func Usecase(ctx context.Context) *UsecaseContainer {
	return &UsecaseContainer{ctx, Repository(ctx)}
}

// CrawlFeeds use case
func (c *UsecaseContainer) CrawlFeeds() *usecase.CrawlFeeds {
	return usecase.NewCrawlFeeds(
		log.NewAppengineLogger(c.ctx),
		c.CrawlFeed(),
	)
}

// CrawlFeed use case
func (c *UsecaseContainer) CrawlFeed() *usecase.CrawlFeed {
	return usecase.NewCrawlFeed(
		log.NewAppengineLogger(c.ctx),
		crawler.New(c.ctx),
		event.NewTaskQueue(c.ctx),
		c.repo.LatestEntryRepository(),
	)
}

// EnqueueTweets use case
func (c *UsecaseContainer) EnqueueTweets() *usecase.EnqueueTweets {
	return usecase.NewEnqueueTweets(
		log.NewAppengineLogger(c.ctx),
		event.NewTaskQueue(c.ctx),
		dao.NewDatastoreTransactor(c.ctx),
		c.repo.TweetItemRepository(),
	)
}

// EnqueueLines use case
func (c *UsecaseContainer) EnqueueLines() *usecase.EnqueueLines {
	return usecase.NewEnqueueLines(
		log.NewAppengineLogger(c.ctx),
		event.NewTaskQueue(c.ctx),
		dao.NewDatastoreTransactor(c.ctx),
		c.repo.LineItemRepository(),
	)
}

// Tweet use case
func (c *UsecaseContainer) Tweet() *usecase.Tweet {
	return usecase.NewTweet(
		log.NewAppengineLogger(c.ctx),
		event.NewTaskQueue(c.ctx),
		twitter.NewTweeter(c.ctx),
	)
}

// Remind use case
func (c *UsecaseContainer) Remind() *usecase.Remind {
	return usecase.NewRemind(
		log.NewAppengineLogger(c.ctx),
		event.NewTaskQueue(c.ctx),
		c.repo.ReminderRepository(),
	)
}

// CheckUstreamStatus use case
func (c *UsecaseContainer) CheckUstreamStatus() *usecase.CheckUstreamStatus {
	return usecase.NewCheckUstreamStatus(
		log.NewAppengineLogger(c.ctx),
		event.NewTaskQueue(c.ctx),
		ustream.NewStatusChecker(c.ctx),
		c.repo.UstreamStatusRepository(),
	)
}

// AddLineNotification use case
func (c *UsecaseContainer) AddLineNotification() *usecase.AddLineNotification {
	return usecase.NewAddLineNotification(
		log.NewAppengineLogger(c.ctx),
		linenotify.NewToken(c.ctx),
		c.repo.LineNotificationRepository(),
	)
}

// HandleLineBotEvents use case
func (c *UsecaseContainer) HandleLineBotEvents() *usecase.HandleLineBotEvents {
	return usecase.NewHandleLineBotEvents(
		log.NewAppengineLogger(c.ctx),
		linebot.New(c.ctx),
		customsearch.NewImageSearcher(c.ctx),
	)
}
