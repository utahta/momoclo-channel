package crawler

import (
	"testing"
	"os"
)

func TestAeNewsChannelParse(t *testing.T) {
	c := NewAeNewsChannel()
	fp, err := os.Open("testdata/ae_news/list_20160715.html")
	if err != nil {
		t.Error("Failed to open ae_news testdata")
	}
	defer fp.Close()

	items, err := c.Parse(fp)
	if err != nil {
		t.Errorf("Failed to ae_news parse. error:%v", err)
	}

	expectedLen := 10
	if len(items) != expectedLen {
		t.Errorf("Invalid item length. %d = %d", len(items), expectedLen)
	}
}
