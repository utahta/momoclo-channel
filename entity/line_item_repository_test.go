package entity

import (
	"testing"
	"time"

	"github.com/utahta/momoclo-channel/crawler"
	"github.com/utahta/momoclo-channel/dao"
	"github.com/utahta/momoclo-channel/testutil"
	"google.golang.org/appengine/aetest"
)

func TestLineItemRepository_Tx(t *testing.T) {
	ctx, done, err := testutil.NewContext(&aetest.Options{StronglyConsistentDatastore: true})
	if err != nil {
		t.Error(err)
	}
	defer done()

	feedItem := crawler.FeedItem{
		EntryTitle:  "entry title",
		EntryURL:    "http://localhost/1",
		PublishedAt: time.Now(),
	}
	item := NewLineItem(
		feedItem.UniqueURL(),
		feedItem.EntryTitle,
		feedItem.EntryURL,
		feedItem.PublishedAt,
		feedItem.ImageURLs,
		feedItem.VideoURLs,
	)
	repo := NewLineItemRepository(dao.NewDatastoreHandler(ctx))
	if err := repo.Save(item); err != nil {
		t.Fatal(err)
	}

	tran := dao.NewDatastoreTransactor(ctx)
	tran.RunInTransaction(func(h dao.PersistenceHandler) error {
		repo := repo.Tx(h)

		item.Title = "entry title z"
		if err := repo.Save(item); err != nil {
			t.Fatal(err)
		}

		v, err := repo.Find(item.ID)
		if err != nil {
			t.Fatal(err)
		}

		// not commit yet
		if v.Title != "entry title" {
			t.Errorf("Expected entry title, got %v", v.Title)
		}
		return nil
	}, nil)

	v, err := repo.Find(item.ID)
	if err != nil {
		t.Fatal(err)
	}

	if v.Title != "entry title z" {
		t.Errorf("Expected title z, got %v", v.Title)
	}
}
