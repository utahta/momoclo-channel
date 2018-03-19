package persistence

import (
	"testing"
	"time"

	"github.com/utahta/momoclo-channel/domain/model"
	"github.com/utahta/momoclo-channel/infrastructure/dao"
	"github.com/utahta/momoclo-channel/testutil"
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
	reminders := []*model.Reminder{
		model.NewReminderOnce("test1", time.Now()),
		model.NewReminderOnce("test2", time.Now()),
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
		t.Errorf("Expected text test1, got %d", reminders[0].Text)
	}
}
