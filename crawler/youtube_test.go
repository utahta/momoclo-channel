package crawler

import (
	"testing"
	"os"
)

func TestYoutubeChannelParse(t *testing.T) {
	c := NewYoutubeChannel()
	fp, err := os.Open("testdata/youtube/feed_20160715.xml")
	if err != nil {
		t.Error("Failed to open youtube testdata")
	}
	defer fp.Close()

	items, err := c.Parse(fp)
	if err != nil {
		t.Errorf("Failed to parse. error:%v", err)
	}

	expectedLen := 15
	if len(items) != expectedLen {
		t.Errorf("Invalid items length. %d = %d", len(items), expectedLen)
	}
}
