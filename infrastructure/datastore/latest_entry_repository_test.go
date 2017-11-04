package datastore

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/utahta/momoclo-channel/domain/entity"
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

	repo := NewLatestEntryRepository()
	for _, test := range tests {
		l, err := entity.ParseLatestEntry(test.url)
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

		if err := repo.Save(ctx, l); err != nil {
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

	repo := NewLatestEntryRepository()
	tests := []struct {
		expectCode string
		expectURL  string
		fn         func(context.Context) string
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
		if err := repo.Save(ctx, blog); err != nil {
			t.Fatal(err)
		}
	}
	time.Sleep(time.Second) // Due to eventual consistency

	for _, test := range tests {
		url := repo.getURL(ctx, test.expectCode)
		if url != test.expectURL {
			t.Fatalf("Expected URL %s, got %s", test.expectURL, url)
		}

		url = test.fn(ctx)
		if url != test.expectURL {
			t.Fatalf("Expected URL %s, got %s", test.expectURL, url)
		}
	}

	url := repo.getURL(ctx, "unknown")
	if url != "" {
		t.Fatalf("Expected URL empty, got %s", url)
	}
}
