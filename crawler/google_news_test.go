package crawler

import (
	"os"
	"testing"
)

func TestGoogleNewsChannelParser(t *testing.T) {
	c := NewGoogleNewsChannelClient()
	fp, err := os.Open("testdata/google_news/feed_20160715.xml")
	if err != nil {
		t.Error("Failed to open youtube testdata")
	}
	defer fp.Close()

	items, err := c.parser.Parse(fp)
	if err != nil {
		t.Errorf("Failed to parse. error:%v", err)
	}

	expectedLen := 10
	if len(items) != expectedLen {
		t.Errorf("Invalid items length. %d = %d", len(items), expectedLen)
	}
}
