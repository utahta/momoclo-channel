package model

import (
	"fmt"
	"testing"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/aetest"
)

func TestPutLatestBlogPost(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Error(err)
	}
	defer done()

	tests := []struct {
		url          string
		expectExists bool
	}{
		{fmt.Sprintf("http://ameblo.jp/%s", BlogPostCodeTamai), true},
		{fmt.Sprintf("http://ameblo.jp/%s", BlogPostCodeMomota), true},
		{fmt.Sprintf("http://ameblo.jp/%s", BlogPostCodeAriyasu), true},
		{fmt.Sprintf("http://ameblo.jp/%s", BlogPostCodeSasaki), true},
		{fmt.Sprintf("http://ameblo.jp/%s", BlogPostCodeTakagi), true},
		{fmt.Sprintf("http://ameblo.jp/%s", "aaa"), false},
	}
	for _, test := range tests {
		l, err := PutLatestBlogPost(ctx, test.url)
		if err != nil {
			t.Fatal(err)
		}

		exists := l != nil
		if exists != test.expectExists {
			t.Fatalf("Expected exists %v, got %v", exists, test.expectExists)
		}
	}
}

func TestGetLatestBlogPostURL(t *testing.T) {
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
		{BlogPostCodeTamai, "http://example.com/1", GetTamaiLatestBlogPostURL},
		{BlogPostCodeMomota, "http://example.com/2", GetMomotaLatestBlogPostURL},
		{BlogPostCodeAriyasu, "http://example.com/3", GetAriyasuLatestBlogPostURL},
		{BlogPostCodeSasaki, "http://example.com/4", GetSasakiLatestBlogPostURL},
		{BlogPostCodeTakagi, "http://example.com/5", GetTakagiLatestBlogPostURL},
	}
	for _, test := range tests {
		blog := NewLatestBlogPost(test.expectCode, test.expectURL)
		if err := blog.Put(ctx); err != nil {
			t.Fatal(err)
		}
	}
	time.Sleep(time.Second) // Due to eventual consistency

	for _, test := range tests {
		url := getLatestBlogPostURL(ctx, test.expectCode)
		if url != test.expectURL {
			t.Fatalf("Expected URL %s, got %s", test.expectURL, url)
		}

		url = test.fn(ctx)
		if url != test.expectURL {
			t.Fatalf("Expected URL %s, got %s", test.expectURL, url)
		}
	}
}
