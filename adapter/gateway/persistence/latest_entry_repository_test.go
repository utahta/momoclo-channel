package persistence_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/utahta/momoclo-channel/adapter/gateway/persistence"
	"github.com/utahta/momoclo-channel/domain/model"
	"github.com/utahta/momoclo-channel/infrastructure/dao"
	"github.com/utahta/momoclo-channel/testutil"
	"google.golang.org/appengine/aetest"
)

func TestLatestEntryRepository_Save(t *testing.T) {
	ctx, done, err := testutil.NewContex(&aetest.Options{StronglyConsistentDatastore: true})
	if err != nil {
		t.Error(err)
	}
	defer done()

	tests := []struct {
		url             string
		expectedSuccess bool
	}{
		{fmt.Sprintf("https://ameblo.jp/%s", model.FeedCodeTamai), true},
		{fmt.Sprintf("https://ameblo.jp/%s", model.FeedCodeMomota), true},
		{fmt.Sprintf("https://ameblo.jp/%s", model.FeedCodeAriyasu), true},
		{fmt.Sprintf("https://ameblo.jp/%s", model.FeedCodeSasaki), true},
		{fmt.Sprintf("https://ameblo.jp/%s", model.FeedCodeTakagi), true},
		{fmt.Sprintf("https://ameblo.jp/%s", "aaa"), false},
		{fmt.Sprintf("http://ameblo.jp/%s", model.FeedCodeMomota), false},
		{"http://www.tfm.co.jp/clover/", true},
	}

	repo := persistence.NewLatestEntryRepository(dao.NewDatastoreHandler(ctx))
	for _, test := range tests {
		l, err := model.NewLatestEntry(test.url)
		if test.expectedSuccess {
			if err != nil {
				t.Fatal(err)
			}
		} else {
			if !strings.Contains(err.Error(), "code not supported") {
				t.Fatal(err)
			}
			continue
		}

		l.PublishedAt = time.Now()
		if err := repo.Save(l); err != nil {
			t.Fatal(err)
		}

		ll, err := repo.FindOrNewByURL(l.URL)
		if err != nil {
			t.Fatal(err)
		}

		if ll.Code == "" {
			t.Error("Expected got code, but empty")
		}
	}

	if err := repo.Save(&model.LatestEntry{ID: "fail-test", URL: "unknown", PublishedAt: time.Now()}); err == nil {
		t.Errorf("Expected got error, but nil")
	}
	if err := repo.Save(&model.LatestEntry{ID: "fail-test", URL: "http://localhost"}); err == nil {
		t.Errorf("Expected got error, but nil")
	}
}

func TestLatestEntryRepository_GetURL(t *testing.T) {
	ctx, done, err := testutil.NewContex(&aetest.Options{StronglyConsistentDatastore: true})
	if err != nil {
		t.Error(err)
	}
	defer done()

	repo := persistence.NewLatestEntryRepository(dao.NewDatastoreHandler(ctx))
	tests := []struct {
		code        model.FeedCode
		expectedURL string
	}{
		{model.FeedCodeTamai, "http://example.com/1"},
		{model.FeedCodeMomota, "http://example.com/2"},
		{model.FeedCodeAriyasu, "http://example.com/3"},
		{model.FeedCodeSasaki, "http://example.com/4"},
		{model.FeedCodeTakagi, "http://example.com/5"},
		{model.FeedCodeHappyclo, "http://example.com/6"},
	}
	for _, test := range tests {
		blog := &model.LatestEntry{ID: test.code.String(), Code: test.code, URL: test.expectedURL, PublishedAt: time.Now()}
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
