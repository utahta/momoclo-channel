package usecase_test

import (
	"testing"

	"github.com/go-playground/validator"
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/adapter/gateway/api/twitter"
	"github.com/utahta/momoclo-channel/container"
	"github.com/utahta/momoclo-channel/domain/model"
	"github.com/utahta/momoclo-channel/event/eventtest"
	"github.com/utahta/momoclo-channel/testutil"
	"github.com/utahta/momoclo-channel/usecase"
	"google.golang.org/appengine/aetest"
)

func TestTweet_Do(t *testing.T) {
	ctx, done, err := testutil.NewContex(&aetest.Options{StronglyConsistentDatastore: true})
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	taskQueue := eventtest.NewTaskQueue()
	u := usecase.NewTweet(container.Logger(ctx).AE(), taskQueue, twitter.NewNopTweeter())

	validationTests := []struct {
		params usecase.TweetParams
	}{
		{},
		{usecase.TweetParams{Requests: []model.TweetRequest{
			{ImageURLs: []string{"a"}},
		}}},
		{usecase.TweetParams{Requests: []model.TweetRequest{
			{VideoURL: "a"},
		}}},
	}

	for _, test := range validationTests {
		err = u.Do(test.params)
		if errs, ok := errors.Cause(err).(validator.ValidationErrors); !ok {
			t.Errorf("Expected validation error, got %v", errs)
		}
	}

	err = u.Do(usecase.TweetParams{Requests: []model.TweetRequest{
		{Text: "test", ImageURLs: []string{"http://localhost/a", "http://localhost/b"}},
		{ImageURLs: []string{"http://localhost/c"}},
		{VideoURL: "http://localhost/d"},
	}})
	if err != nil {
		t.Fatal(err)
	}

	if len(taskQueue.Tasks) != 1 {
		t.Errorf("Expected taskqueue length 1, got %v", len(taskQueue.Tasks))
	}
}
