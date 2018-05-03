package entity

import (
	"fmt"
	"testing"
	"time"

	"github.com/utahta/momoclo-channel/crawler"
	"github.com/utahta/momoclo-channel/dao"
	"github.com/utahta/momoclo-channel/testutil"
	"google.golang.org/appengine/aetest"
)

func TestLatestEntryRepository_Save(t *testing.T) {
	ctx, done, err := testutil.NewContext(&aetest.Options{StronglyConsistentDatastore: true})
	if err != nil {
		t.Error(err)
	}
	defer done()

	tests := []struct {
		code crawler.FeedCode
		url  string
	}{
		{crawler.FeedCodeTamai, fmt.Sprintf("https://ameblo.jp/%s", crawler.FeedCodeTamai)},
		{crawler.FeedCodeMomota, fmt.Sprintf("https://ameblo.jp/%s", crawler.FeedCodeMomota)},
		{crawler.FeedCodeSasaki, fmt.Sprintf("https://ameblo.jp/%s", crawler.FeedCodeSasaki)},
		{crawler.FeedCodeTakagi, fmt.Sprintf("https://ameblo.jp/%s", crawler.FeedCodeTakagi)},
		{crawler.FeedCodeHappyclo, "http://www.tfm.co.jp/clover/"},
	}

	repo := NewLatestEntryRepository(dao.NewDatastoreHandler(ctx))
	for _, test := range tests {
		l, err := NewLatestEntry(test.code.String(), test.url)
		if err != nil {
			t.Fatal(err)
		}

		l.PublishedAt = time.Now()
		if err := repo.Save(l); err != nil {
			t.Fatal(err)
		}

		ll, err := repo.FindOrNewByURL(l.Code, l.URL)
		if err != nil {
			t.Fatal(err)
		}

		if ll.Code == "" {
			t.Error("Expected got code, but empty")
		}
	}

	if err := repo.Save(&LatestEntry{ID: "fail-test", URL: "unknown", PublishedAt: time.Now()}); err == nil {
		t.Errorf("Expected got error, but nil")
	}
	if err := repo.Save(&LatestEntry{ID: "fail-test", URL: "http://localhost"}); err == nil {
		t.Errorf("Expected got error, but nil")
	}
}

func TestLatestEntryRepository_GetURL(t *testing.T) {
	ctx, done, err := testutil.NewContext(&aetest.Options{StronglyConsistentDatastore: true})
	if err != nil {
		t.Error(err)
	}
	defer done()

	repo := NewLatestEntryRepository(dao.NewDatastoreHandler(ctx))
	tests := []struct {
		code        crawler.FeedCode
		expectedURL string
	}{
		{crawler.FeedCodeTamai, "http://example.com/1"},
		{crawler.FeedCodeMomota, "http://example.com/2"},
		{crawler.FeedCodeSasaki, "http://example.com/4"},
		{crawler.FeedCodeTakagi, "http://example.com/5"},
		{crawler.FeedCodeHappyclo, "http://example.com/6"},
	}
	for _, test := range tests {
		blog := &LatestEntry{ID: test.code.String(), Code: test.code.String(), URL: test.expectedURL, PublishedAt: time.Now()}
		if err := repo.Save(blog); err != nil {
			t.Fatal(err)
		}
	}

	for _, test := range tests {
		url := repo.GetURL(test.code.String())
		if url != test.expectedURL {
			t.Fatalf("Expected URL %s, got %s", test.expectedURL, url)
		}
	}
}
