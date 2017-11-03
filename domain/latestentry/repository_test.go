package latestentry

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/utahta/momoclo-channel/domain"
	"google.golang.org/appengine/aetest"
)

func TestPutLatestEntry(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Error(err)
	}
	defer done()

	tests := []struct {
		url          string
		expectExists bool
	}{
		{fmt.Sprintf("https://ameblo.jp/%s", domain.LatestEntryCodeTamai), true},
		{fmt.Sprintf("https://ameblo.jp/%s", domain.LatestEntryCodeMomota), true},
		{fmt.Sprintf("https://ameblo.jp/%s", domain.LatestEntryCodeAriyasu), true},
		{fmt.Sprintf("https://ameblo.jp/%s", domain.LatestEntryCodeSasaki), true},
		{fmt.Sprintf("https://ameblo.jp/%s", domain.LatestEntryCodeTakagi), true},
		{fmt.Sprintf("https://ameblo.jp/%s", "aaa"), false},
		{fmt.Sprintf("http://ameblo.jp/%s", domain.LatestEntryCodeMomota), false},
		{"http://www.tfm.co.jp/clover/", true},
	}
	for _, test := range tests {
		l, err := Repository.PutURL(ctx, test.url)
		if err != nil {
			t.Fatal(err)
		}

		exists := l != nil
		if exists != test.expectExists {
			t.Fatalf("Expected exists %v, got %v", exists, test.expectExists)
		}
	}
}

func TestGetLatestEntryURL(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Error(err)
	}
	defer done()

	tests := []struct {
		expectCode string
		expectURL  string
		fn         func(context.Context) string
	}{
		{domain.LatestEntryCodeTamai, "http://example.com/1", Repository.GetTamaiURL},
		{domain.LatestEntryCodeMomota, "http://example.com/2", Repository.GetMomotaURL},
		{domain.LatestEntryCodeAriyasu, "http://example.com/3", Repository.GetAriyasuURL},
		{domain.LatestEntryCodeSasaki, "http://example.com/4", Repository.GetSasakiURL},
		{domain.LatestEntryCodeTakagi, "http://example.com/5", Repository.GetTakagiURL},
		{domain.LatestEntryCodeHappyclo, "http://example.com/6", Repository.GetHappycloURL},
	}
	for _, test := range tests {
		blog := domain.NewLatestEntry(test.expectCode, test.expectURL)
		if err := blog.Put(ctx); err != nil {
			t.Fatal(err)
		}
	}
	time.Sleep(time.Second) // Due to eventual consistency

	for _, test := range tests {
		url := Repository.getURL(ctx, test.expectCode)
		if url != test.expectURL {
			t.Fatalf("Expected URL %s, got %s", test.expectURL, url)
		}

		url = test.fn(ctx)
		if url != test.expectURL {
			t.Fatalf("Expected URL %s, got %s", test.expectURL, url)
		}
	}

	url := Repository.getURL(ctx, "unknown")
	if url != "" {
		t.Fatalf("Expected URL empty, got %s", url)
	}
}
