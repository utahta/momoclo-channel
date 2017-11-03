package reminder

import (
	"testing"
	"time"

	"github.com/utahta/momoclo-channel/domain"
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
)

func TestReminderQuery_GetAll(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	now := time.Now()
	res := []*domain.Reminder{domain.NewReminderOnce("test1", now), domain.NewReminderOnce("test2", now)}
	res[1].Enabled = false
	for _, re := range res {
		if err := re.Put(ctx); err != nil {
			t.Fatal(err)
		}
	}
	time.Sleep(time.Second) // Due to eventual consistency

	num, err := datastore.NewQuery("Reminder").Count(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if num != 2 {
		t.Fatalf("Expected len 2, got %d", len(res))
	}

	res, err = Repository.GetAll(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 {
		t.Fatalf("Expected len 1, got %d", len(res))
	}
	if res[0].Text != "test1" {
		t.Errorf("Expected text test1, got %d", res[0].Text)
	}
}
