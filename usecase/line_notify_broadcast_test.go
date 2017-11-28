package usecase_test

import (
	"testing"

	"fmt"

	"github.com/go-playground/validator"
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/container"
	"github.com/utahta/momoclo-channel/domain/model"
	"github.com/utahta/momoclo-channel/infrastructure/event/eventtest"
	"github.com/utahta/momoclo-channel/lib/aetestutil"
	"github.com/utahta/momoclo-channel/lib/config"
	"github.com/utahta/momoclo-channel/usecase"
	"google.golang.org/appengine/aetest"
)

func TestLineNotifyBroadcast_Do(t *testing.T) {
	ctx, done, err := aetestutil.NewContex(&aetest.Options{StronglyConsistentDatastore: true})
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	taskQueue := eventtest.NewTaskQueue()
	repo := container.Repository(ctx).LineNotificationRepository()
	u := usecase.NewLineNotifyBroadcast(container.Logger(ctx).AE(), taskQueue, repo)

	validationTests := []struct {
		params usecase.LineNotifyBroadcastParams
	}{
		{usecase.LineNotifyBroadcastParams{Messages: nil}},
		{usecase.LineNotifyBroadcastParams{Messages: []model.LineNotifyMessage{
			{Text: ""},
		}}},
		{usecase.LineNotifyBroadcastParams{Messages: []model.LineNotifyMessage{
			{Text: "hello", ImageURL: "unknown"},
		}}},
	}

	for _, test := range validationTests {
		err = u.Do(test.params)
		if errs, ok := errors.Cause(err).(validator.ValidationErrors); !ok {
			t.Errorf("Expected validation error, got %v. params:%v", errs, test.params)
		}
	}

	config.C = config.Config{LineNotify: config.LineNotify{TokenKey: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}}
	for i := 0; i < 10; i++ {
		l, err := model.NewLineNotification(config.C.LineNotify.TokenKey, fmt.Sprintf("token-%v", i))
		if err != nil {
			t.Fatal(err)
		}
		repo.Save(l)
	}

	err = u.Do(usecase.LineNotifyBroadcastParams{Messages: []model.LineNotifyMessage{
		{Text: "hello"},
		{Text: " ", ImageURL: "http://localhost/a"},
	},
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(taskQueue.Tasks) != 10 {
		t.Errorf("Expected taskqueue length 10, got %v", len(taskQueue.Tasks))
	}
}