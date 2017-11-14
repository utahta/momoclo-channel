package persistence_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/utahta/momoclo-channel/adapter/gateway/persistence"
	"github.com/utahta/momoclo-channel/domain/model"
	"github.com/utahta/momoclo-channel/domain/service/latestentry"
	"github.com/utahta/momoclo-channel/infrastructure/dao"
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
		{fmt.Sprintf("https://ameblo.jp/%s", model.LatestEntryCodeTamai), true},
		{fmt.Sprintf("https://ameblo.jp/%s", model.LatestEntryCodeMomota), true},
		{fmt.Sprintf("https://ameblo.jp/%s", model.LatestEntryCodeAriyasu), true},
		{fmt.Sprintf("https://ameblo.jp/%s", model.LatestEntryCodeSasaki), true},
		{fmt.Sprintf("https://ameblo.jp/%s", model.LatestEntryCodeTakagi), true},
		{fmt.Sprintf("https://ameblo.jp/%s", "aaa"), false},
		{fmt.Sprintf("http://ameblo.jp/%s", model.LatestEntryCodeMomota), false},
		{"http://www.tfm.co.jp/clover/", true},
	}

	repo := persistence.NewLatestEntryRepository(dao.NewDatastoreHandler(ctx))
	for _, test := range tests {
		l, err := latestentry.Parse(test.url)
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

	repo := persistence.NewLatestEntryRepository(dao.NewDatastoreHandler(ctx))
	tests := []struct {
		code        string
		expectedURL string
	}{
		{model.LatestEntryCodeTamai, "http://example.com/1"},
		{model.LatestEntryCodeMomota, "http://example.com/2"},
		{model.LatestEntryCodeAriyasu, "http://example.com/3"},
		{model.LatestEntryCodeSasaki, "http://example.com/4"},
		{model.LatestEntryCodeTakagi, "http://example.com/5"},
		{model.LatestEntryCodeHappyclo, "http://example.com/6"},
	}
	for _, test := range tests {
		blog := &model.LatestEntry{ID: test.code, Code: test.code, URL: test.expectedURL}
		if err := repo.Save(blog); err != nil {
			t.Fatal(err)
		}
	}
	time.Sleep(time.Second) // Due to eventual consistency

	for _, test := range tests {
		url := repo.GetURL(test.code)
		if url != test.expectedURL {
			t.Fatalf("Expected URL %s, got %s", test.expectedURL, url)
		}
	}
}
