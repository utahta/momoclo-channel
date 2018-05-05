package container

import (
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
	repo   *RepositoryContainer
	logger *LoggerContainer
}

// Usecase returns container of use case
func Usecase() *UsecaseContainer {
	return &UsecaseContainer{Repository(), Logger()}
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
		crawler.New(),
		event.NewTaskQueue(),
		c.repo.LatestEntryRepository(),
	)
}

// EnqueueTweets use case
func (c *UsecaseContainer) EnqueueTweets() *usecase.EnqueueTweets {
	return usecase.NewEnqueueTweets(
		c.logger.AE(),
		event.NewTaskQueue(),
		dao.NewDatastoreTransactor(),
		c.repo.TweetItemRepository(),
	)
}

// EnqueueLines use case
func (c *UsecaseContainer) EnqueueLines() *usecase.EnqueueLines {
	return usecase.NewEnqueueLines(
		c.logger.AE(),
		event.NewTaskQueue(),
		dao.NewDatastoreTransactor(),
		c.repo.LineItemRepository(),
	)
}

// Tweet use case
func (c *UsecaseContainer) Tweet() *usecase.Tweet {
	return usecase.NewTweet(
		c.logger.AE(),
		event.NewTaskQueue(),
		twitter.NewTweeter(),
	)
}

// Remind use case
func (c *UsecaseContainer) Remind() *usecase.Remind {
	return usecase.NewRemind(
		c.logger.AE(),
		event.NewTaskQueue(),
		c.repo.ReminderRepository(),
	)
}

// CheckUstream use case
func (c *UsecaseContainer) CheckUstream() *usecase.CheckUstream {
	return usecase.NewCheckUstream(
		c.logger.AE(),
		event.NewTaskQueue(),
		ustream.NewStatusChecker(),
		c.repo.UstreamStatusRepository(),
	)
}

// AddLineNotification use case
func (c *UsecaseContainer) AddLineNotification() *usecase.AddLineNotification {
	return usecase.NewAddLineNotification(
		c.logger.AE(),
		linenotify.NewToken(),
		c.repo.LineNotificationRepository(),
	)
}

// HandleLineBotEvents use case
func (c *UsecaseContainer) HandleLineBotEvents() *usecase.HandleLineBotEvents {
	return usecase.NewHandleLineBotEvents(
		c.logger.AE(),
		linebot.New(),
		customsearch.NewImageSearcher(),
	)
}

// LineNotifyBroadcast use case
func (c *UsecaseContainer) LineNotifyBroadcast() *usecase.LineNotifyBroadcast {
	return usecase.NewLineNotifyBroadcast(
		c.logger.AE(),
		event.NewTaskQueue(),
		c.repo.LineNotificationRepository(),
	)
}

// LineNotify use case
func (c *UsecaseContainer) LineNotify() *usecase.LineNotify {
	return usecase.NewLineNotify(
		c.logger.AE(),
		event.NewTaskQueue(),
		linenotify.New(),
		c.repo.LineNotificationRepository(),
	)
}
