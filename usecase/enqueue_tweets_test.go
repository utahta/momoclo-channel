package usecase_test

import (
	"testing"
	"time"

	"github.com/go-playground/validator"
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/container"
	"github.com/utahta/momoclo-channel/dao"
	"github.com/utahta/momoclo-channel/event/eventtest"
	"github.com/utahta/momoclo-channel/testutil"
	"github.com/utahta/momoclo-channel/types"
	"github.com/utahta/momoclo-channel/usecase"
	"google.golang.org/appengine/aetest"
)

func TestEnqueueTweets_Do(t *testing.T) {
	ctx, done, err := testutil.NewContex(&aetest.Options{StronglyConsistentDatastore: true})
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	taskQueue := eventtest.NewTaskQueue()
	repo := container.Repository(ctx).TweetItemRepository()
	u := usecase.NewEnqueueTweets(container.Logger(ctx).AE(), taskQueue, dao.NewDatastoreTransactor(ctx), repo)
	publishedAt, _ := time.Parse("2006-01-02 15:04:05", "2008-05-17 00:00:00")
	feedItem := types.FeedItem{
		Title:       "title",
		URL:         "http://localhost",
		EntryTitle:  "entry_title",
		EntryURL:    "http://localhost/entry",
		ImageURLs:   []string{"http://localhost/img_1", "http://localhost/img_2"},
		VideoURLs:   []string{"http://localhost/mp4_1"},
		PublishedAt: publishedAt,
	}
	item := types.NewTweetItem(feedItem)
	if repo.Exists(item.ID) {
		t.Errorf("Expected tweet item not found, but exists. feedItem:%v", feedItem)
	}

	err = u.Do(usecase.EnqueueTweetsParams{FeedItem: types.FeedItem{}})
	if errs, ok := errors.Cause(err).(validator.ValidationErrors); !ok {
		t.Errorf("Expected validation errors, got %v", errs)
	}

	if err := u.Do(usecase.EnqueueTweetsParams{FeedItem: feedItem}); err != nil {
		t.Fatal(err)
	}

	if !repo.Exists(item.ID) {
		t.Errorf("Expected tweet item exists, but not found. feedItem:%v", feedItem)
	}

	if len(taskQueue.Tasks) != 1 {
		t.Errorf("Expected task length 1, got %v", len(taskQueue.Tasks))
	}
	if taskQueue.Tasks[0].QueueName != "queue-tweet" {
		t.Errorf("Expected queue name queue-tweet, got %v", taskQueue.Tasks[0].QueueName)
	}
	if taskQueue.Tasks[0].Path != "/tweet" {
		t.Errorf("Expected queue path /queue/tweet, got %v", taskQueue.Tasks[0].Path)
	}
}
