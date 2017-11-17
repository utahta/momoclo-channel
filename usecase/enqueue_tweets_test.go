package usecase_test

import (
	"testing"
	"time"

	"github.com/utahta/momoclo-channel/container"
	"github.com/utahta/momoclo-channel/domain/model"
	"github.com/utahta/momoclo-channel/domain/service/convert"
	"github.com/utahta/momoclo-channel/infrastructure/dao"
	"github.com/utahta/momoclo-channel/infrastructure/event/eventtest"
	"github.com/utahta/momoclo-channel/infrastructure/log"
	"github.com/utahta/momoclo-channel/lib/aetestutil"
	"github.com/utahta/momoclo-channel/usecase"
	"google.golang.org/appengine/aetest"
)

func TestEnqueueTweets_Do(t *testing.T) {
	ctx, done, err := aetestutil.NewContex(&aetest.Options{StronglyConsistentDatastore: true})
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	logger := log.NewAppengineLogger(ctx)
	taskQueue := eventtest.NewTaskQueue(logger)
	repo := container.Repository(ctx).TweetItemRepository()
	u := usecase.NewEnqueueTweets(logger, taskQueue, dao.NewDatastoreTransactor(ctx), repo)
	publishedAt, _ := time.Parse("2006-01-02 15:04:05", "2008-05-17 00:00:00")
	feedItem := model.FeedItem{
		Title:       "title",
		URL:         "http://localhost",
		EntryTitle:  "entry_title",
		EntryURL:    "http://localhost/entry",
		ImageURLs:   []string{"http://localhost/img_1", "http://localhost/img_2"},
		VideoURLs:   []string{"http://localhost/mp4_1"},
		PublishedAt: publishedAt,
	}
	tweetItem := convert.FeedItemToTweetItem(feedItem)
	if repo.Exists(tweetItem.ID) {
		t.Errorf("Expected tweet item not found, but exists. feedItem:%v", feedItem)
	}

	if err := u.Do(usecase.EnqueueTweetsParams{FeedItem: feedItem}); err != nil {
		t.Fatal(err)
	}

	if !repo.Exists(tweetItem.ID) {
		t.Errorf("Expected tweet item exists, but not found. feedItem:%v", feedItem)
	}

	if len(taskQueue.Tasks) != 1 {
		t.Errorf("Expected task length 1, got %v", len(taskQueue.Tasks))
	}
	if taskQueue.Tasks[0].QueueName != "queue-tweet" {
		t.Errorf("Expected queue name queue-tweet, got %v", taskQueue.Tasks[0].QueueName)
	}
	if taskQueue.Tasks[0].Path != "/queue/tweet" {
		t.Errorf("Expected queue path /queue/tweet, got %v", taskQueue.Tasks[0].Path)
	}
}
