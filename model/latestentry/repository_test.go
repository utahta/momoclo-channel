package latestentry

import (
	"fmt"
	"testing"
	"time"

	"github.com/utahta/momoclo-channel/model"
	"golang.org/x/net/context"
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
		{fmt.Sprintf("https://ameblo.jp/%s", model.LatestEntryCodeTamai), true},
		{fmt.Sprintf("https://ameblo.jp/%s", model.LatestEntryCodeMomota), true},
		{fmt.Sprintf("https://ameblo.jp/%s", model.LatestEntryCodeAriyasu), true},
		{fmt.Sprintf("https://ameblo.jp/%s", model.LatestEntryCodeSasaki), true},
		{fmt.Sprintf("https://ameblo.jp/%s", model.LatestEntryCodeTakagi), true},
		{fmt.Sprintf("https://ameblo.jp/%s", "aaa"), false},
		{fmt.Sprintf("http://ameblo.jp/%s", model.LatestEntryCodeMomota), false},
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
		{model.LatestEntryCodeTamai, "http://example.com/1", Repository.GetTamaiURL},
		{model.LatestEntryCodeMomota, "http://example.com/2", Repository.GetMomotaURL},
		{model.LatestEntryCodeAriyasu, "http://example.com/3", Repository.GetAriyasuURL},
		{model.LatestEntryCodeSasaki, "http://example.com/4", Repository.GetSasakiURL},
		{model.LatestEntryCodeTakagi, "http://example.com/5", Repository.GetTakagiURL},
		{model.LatestEntryCodeHappyclo, "http://example.com/6", Repository.GetHappycloURL},
	}
	for _, test := range tests {
		blog := model.NewLatestEntry(test.expectCode, test.expectURL)
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
