package model

import (
	"fmt"
	"testing"
	"time"

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
		{fmt.Sprintf("http://ameblo.jp/%s", LatestEntryCodeTamai), true},
		{fmt.Sprintf("http://ameblo.jp/%s", LatestEntryCodeMomota), true},
		{fmt.Sprintf("http://ameblo.jp/%s", LatestEntryCodeAriyasu), true},
		{fmt.Sprintf("http://ameblo.jp/%s", LatestEntryCodeSasaki), true},
		{fmt.Sprintf("http://ameblo.jp/%s", LatestEntryCodeTakagi), true},
		{fmt.Sprintf("http://ameblo.jp/%s", "aaa"), false},
		{"http://www.tfm.co.jp/clover/", true},
	}
	for _, test := range tests {
		l, err := PutLatestEntry(ctx, test.url)
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
		{LatestEntryCodeTamai, "http://example.com/1", GetTamaiLatestEntryURL},
		{LatestEntryCodeMomota, "http://example.com/2", GetMomotaLatestEntryURL},
		{LatestEntryCodeAriyasu, "http://example.com/3", GetAriyasuLatestEntryURL},
		{LatestEntryCodeSasaki, "http://example.com/4", GetSasakiLatestEntryURL},
		{LatestEntryCodeTakagi, "http://example.com/5", GetTakagiLatestEntryURL},
		{LatestEntryCodeHappyclo, "http://example.com/6", GetHappycloLatestEntryURL},
	}
	for _, test := range tests {
		blog := NewLatestEntry(test.expectCode, test.expectURL)
		if err := blog.Put(ctx); err != nil {
			t.Fatal(err)
		}
	}
	time.Sleep(time.Second) // Due to eventual consistency

	for _, test := range tests {
		url := getLatestEntryURL(ctx, test.expectCode)
		if url != test.expectURL {
			t.Fatalf("Expected URL %s, got %s", test.expectURL, url)
		}

		url = test.fn(ctx)
		if url != test.expectURL {
			t.Fatalf("Expected URL %s, got %s", test.expectURL, url)
		}
	}

	url := getLatestEntryURL(ctx, "unknown")
	if url != "" {
		t.Fatalf("Expected URL empty, got %s", url)
	}
}
