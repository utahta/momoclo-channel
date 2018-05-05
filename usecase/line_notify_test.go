package usecase_test

import (
	"testing"

	"github.com/go-playground/validator"
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/container"
	"github.com/utahta/momoclo-channel/event/eventtest"
	"github.com/utahta/momoclo-channel/linenotify"
	"github.com/utahta/momoclo-channel/testutil"
	"github.com/utahta/momoclo-channel/usecase"
	"google.golang.org/appengine/aetest"
)

func TestLineNotify_Do(t *testing.T) {
	ctx, done, err := testutil.NewContext(&aetest.Options{StronglyConsistentDatastore: true})
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	taskQueue := eventtest.NewTaskQueue()
	repo := container.Repository().LineNotificationRepository()
	u := usecase.NewLineNotify(container.Logger().AE(), taskQueue, linenotify.NewNop(), repo)

	validationTests := []struct {
		params usecase.LineNotifyParams
	}{
		{usecase.LineNotifyParams{Request: linenotify.Request{ID: "id-1"}}},
		{usecase.LineNotifyParams{Request: linenotify.Request{AccessToken: "token"}}},
		{usecase.LineNotifyParams{Request: linenotify.Request{
			ID: "id-2", AccessToken: "token",
		}}},
		{usecase.LineNotifyParams{Request: linenotify.Request{
			ID: "id-3", AccessToken: "token", Messages: []linenotify.Message{
				{Text: ""},
			},
		}}},
		{usecase.LineNotifyParams{Request: linenotify.Request{
			ID: "id-4", AccessToken: "token", Messages: []linenotify.Message{
				{Text: "hello", ImageURL: "unknown"},
			},
		}}},
	}

	for _, test := range validationTests {
		err = u.Do(ctx, test.params)
		if errs, ok := errors.Cause(err).(validator.ValidationErrors); !ok {
			t.Errorf("Expected validation error, got %v. params:%v", errs, test.params)
		}
	}

	err = u.Do(ctx, usecase.LineNotifyParams{Request: linenotify.Request{
		ID: "id-1", AccessToken: "token", Messages: []linenotify.Message{
			{Text: "hello"},
			{Text: " ", ImageURL: "http://localhost/a"},
		},
	}})
	if err != nil {
		t.Fatal(err)
	}

	if len(taskQueue.Tasks) != 1 {
		t.Errorf("Expected taskqueue length 1, got %v", len(taskQueue.Tasks))
	}
}
