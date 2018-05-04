package container

import (
	"context"

	"github.com/utahta/momoclo-channel/crawler"
	"github.com/utahta/momoclo-channel/customsearch"
	"github.com/utahta/momoclo-channel/dao"
	"github.com/utahta/momoclo-channel/event"
	"github.com/utahta/momoclo-channel/linebot"
	"github.com/utahta/momoclo-channel/linenotify"
	"github.com/utahta/momoclo-channel/twitter"
	"github.com/utahta/momoclo-channel/usecase"
	"github.com/utahta/momoclo-channel/ustream"
)

// UsecaseContainer dependency injection
type UsecaseContainer struct {
	ctx    context.Context
	repo   *RepositoryContainer
	logger *LoggerContainer
}

// Usecase returns container of use case
func Usecase(ctx context.Context) *UsecaseContainer {
	return &UsecaseContainer{ctx, Repository(ctx), Logger(ctx)}
}

// CrawlFeeds use case
func (c *UsecaseContainer) CrawlFeeds() *usecase.CrawlFeeds {
	return usecase.NewCrawlFeeds(
		c.logger.AE(),
		c.CrawlFeed(),
	)
}

// CrawlFeed use case
func (c *UsecaseContainer) CrawlFeed() *usecase.CrawlFeed {
	return usecase.NewCrawlFeed(
		c.logger.AE(),
		crawler.New(c.ctx),
		event.NewTaskQueue(c.ctx),
		c.repo.LatestEntryRepository(),
	)
}

// EnqueueTweets use case
func (c *UsecaseContainer) EnqueueTweets() *usecase.EnqueueTweets {
	return usecase.NewEnqueueTweets(
		c.logger.AE(),
		event.NewTaskQueue(c.ctx),
		dao.NewDatastoreTransactor(c.ctx),
		c.repo.TweetItemRepository(),
	)
}

// EnqueueLines use case
func (c *UsecaseContainer) EnqueueLines() *usecase.EnqueueLines {
	return usecase.NewEnqueueLines(
		c.logger.AE(),
		event.NewTaskQueue(c.ctx),
		dao.NewDatastoreTransactor(c.ctx),
		c.repo.LineItemRepository(),
	)
}

// Tweet use case
func (c *UsecaseContainer) Tweet() *usecase.Tweet {
	return usecase.NewTweet(
		c.logger.AE(),
		event.NewTaskQueue(c.ctx),
		twitter.NewTweeter(c.ctx),
	)
}

// Remind use case
func (c *UsecaseContainer) Remind() *usecase.Remind {
	return usecase.NewRemind(
		c.logger.AE(),
		event.NewTaskQueue(c.ctx),
		c.repo.ReminderRepository(),
	)
}

// CheckUstream use case
func (c *UsecaseContainer) CheckUstream() *usecase.CheckUstream {
	return usecase.NewCheckUstream(
		c.logger.AE(),
		event.NewTaskQueue(c.ctx),
		ustream.NewStatusChecker(c.ctx),
		c.repo.UstreamStatusRepository(),
	)
}

// AddLineNotification use case
func (c *UsecaseContainer) AddLineNotification() *usecase.AddLineNotification {
	return usecase.NewAddLineNotification(
		c.logger.AE(),
		linenotify.NewToken(c.ctx),
		c.repo.LineNotificationRepository(),
	)
}

// HandleLineBotEvents use case
func (c *UsecaseContainer) HandleLineBotEvents() *usecase.HandleLineBotEvents {
	return usecase.NewHandleLineBotEvents(
		c.logger.AE(),
		linebot.New(c.ctx),
		customsearch.MustNewImageSearcher(c.ctx),
	)
}

// LineNotifyBroadcast use case
func (c *UsecaseContainer) LineNotifyBroadcast() *usecase.LineNotifyBroadcast {
	return usecase.NewLineNotifyBroadcast(
		c.logger.AE(),
		event.NewTaskQueue(c.ctx),
		c.repo.LineNotificationRepository(),
	)
}

// LineNotify use case
func (c *UsecaseContainer) LineNotify() *usecase.LineNotify {
	return usecase.NewLineNotify(
		c.logger.AE(),
		event.NewTaskQueue(c.ctx),
		linenotify.New(c.ctx),
		c.repo.LineNotificationRepository(),
	)
}
