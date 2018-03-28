package persistence

import (
	"testing"
	"time"

	"github.com/utahta/momoclo-channel/dao"
	"github.com/utahta/momoclo-channel/testutil"
	"github.com/utahta/momoclo-channel/types"
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
)

func TestReminderRepository_FindAll(t *testing.T) {
	ctx, done, err := testutil.NewContex(&aetest.Options{StronglyConsistentDatastore: true})
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	repo := NewReminderRepository(dao.NewDatastoreHandler(ctx))
	reminders := []*types.Reminder{
		types.NewReminderOnce("test1", time.Now()),
		types.NewReminderOnce("test2", time.Now()),
	}
	reminders[1].Enabled = false
	for _, reminder := range reminders {
		if err := repo.Save(reminder); err != nil {
			t.Fatal(err)
		}
	}

	num, err := datastore.NewQuery("Reminder").Count(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if num != 2 {
		t.Fatalf("Expected len 2, got %d", len(reminders))
	}

	reminders, err = repo.FindAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(reminders) != 1 {
		t.Fatalf("Expected len 1, got %d", len(reminders))
	}
	if reminders[0].Text != "test1" {
		t.Errorf("Expected text test1, got %v", reminders[0].Text)
	}
}
