package persistence_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/utahta/momoclo-channel/adapter/persistence"
	"github.com/utahta/momoclo-channel/domain/entity"
	"github.com/utahta/momoclo-channel/domain/service/latest_entry"
	"github.com/utahta/momoclo-channel/infrastructure/datastore"
	"google.golang.org/appengine/aetest"
)

func TestLatestEntryRepository_Save(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Error(err)
	}
	defer done()

	tests := []struct {
		url             string
		expectedSuccess bool
	}{
		{fmt.Sprintf("https://ameblo.jp/%s", entity.LatestEntryCodeTamai), true},
		{fmt.Sprintf("https://ameblo.jp/%s", entity.LatestEntryCodeMomota), true},
		{fmt.Sprintf("https://ameblo.jp/%s", entity.LatestEntryCodeAriyasu), true},
		{fmt.Sprintf("https://ameblo.jp/%s", entity.LatestEntryCodeSasaki), true},
		{fmt.Sprintf("https://ameblo.jp/%s", entity.LatestEntryCodeTakagi), true},
		{fmt.Sprintf("https://ameblo.jp/%s", "aaa"), false},
		{fmt.Sprintf("http://ameblo.jp/%s", entity.LatestEntryCodeMomota), false},
		{"http://www.tfm.co.jp/clover/", true},
	}

	repo := persistence.NewLatestEntryRepository(datastore.New(ctx))
	for _, test := range tests {
		l, err := latest_entry.Parse(test.url)
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

		if err := repo.Save(l); err != nil {
			t.Fatal(err)
		}
	}
}

func TestLatestEntryRepository_GetURL(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Error(err)
	}
	defer done()

	repo := persistence.NewLatestEntryRepository(datastore.New(ctx))
	tests := []struct {
		expectCode string
		expectURL  string
		fn         func() string
	}{
		{entity.LatestEntryCodeTamai, "http://example.com/1", repo.GetTamaiURL},
		{entity.LatestEntryCodeMomota, "http://example.com/2", repo.GetMomotaURL},
		{entity.LatestEntryCodeAriyasu, "http://example.com/3", repo.GetAriyasuURL},
		{entity.LatestEntryCodeSasaki, "http://example.com/4", repo.GetSasakiURL},
		{entity.LatestEntryCodeTakagi, "http://example.com/5", repo.GetTakagiURL},
		{entity.LatestEntryCodeHappyclo, "http://example.com/6", repo.GetHappycloURL},
	}
	for _, test := range tests {
		blog := &entity.LatestEntry{ID: test.expectCode, Code: test.expectCode, URL: test.expectURL}
		if err := repo.Save(blog); err != nil {
			t.Fatal(err)
		}
	}
	time.Sleep(time.Second) // Due to eventual consistency

	for _, test := range tests {
		url := test.fn()
		if url != test.expectURL {
			t.Fatalf("Expected URL %s, got %s", test.expectURL, url)
		}
	}
}
