package model

import (
	"testing"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/aetest"
)

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
