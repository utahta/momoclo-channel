package crawler

import (
	"os"
	"testing"
)

func TestBlogChannelParser(t *testing.T) {
	var tests = []struct {
		c *ChannelClient
	}{
		{NewTamaiBlogChannelClient(nil)},
		{NewMomotaBlogChannelClient(nil)},
		{NewAriyasuBlogChannelClient(nil)},
		{NewSasakiBlogChannelClient(nil)},
		{NewTakagiBlogChannelClient(nil)},
	}

	for _, test := range tests {
		if test.c.parser == nil {
			t.Errorf("Channel.Parse() is nil. %#v", test.c)
		}
	}
}

func TestBlogChannelParserList(t *testing.T) {
	var tests = []struct {
		c     *ChannelClient
		input string
	}{
		{NewTamaiBlogChannelClient(nil), "testdata/blog/list_tamai_20160714.html"},
		{NewMomotaBlogChannelClient(nil), "testdata/blog/list_momota_20160714.html"},
		{NewAriyasuBlogChannelClient(nil), "testdata/blog/list_ariyasu_20160714.html"},
		{NewSasakiBlogChannelClient(nil), "testdata/blog/list_sasaki_20160714.html"},
		{NewTakagiBlogChannelClient(nil), "testdata/blog/list_takagi_20160714.html"},
	}

	for _, test := range tests {
		func() {
			fp, err := os.Open(test.input)
			if err != nil {
				t.Errorf("Failed to open test list data. input:%s", test.input)
			}
			defer fp.Close()

			parser := test.c.parser.(*blogChannelParser)
			items, err := parser.parseList(fp)
			if err != nil {
				t.Errorf("Failed to parse list. error:%v", err)
			}

			if len(items) != 5 {
				t.Errorf("Invalid items size. len:%d", len(items))
			}
		}()
	}
}

func TestBlogChannelParserListWithOption(t *testing.T) {
	const expectMaxItemNum = 10
	opt := &BlogChannelParserOption{MaxItemNum: expectMaxItemNum}
	var tests = []struct {
		c     *ChannelClient
		input string
	}{
		{NewTamaiBlogChannelClient(opt), "testdata/blog/list_tamai_20160714.html"},
		{NewMomotaBlogChannelClient(opt), "testdata/blog/list_momota_20160714.html"},
		{NewAriyasuBlogChannelClient(opt), "testdata/blog/list_ariyasu_20160714.html"},
		{NewSasakiBlogChannelClient(opt), "testdata/blog/list_sasaki_20160714.html"},
		{NewTakagiBlogChannelClient(opt), "testdata/blog/list_takagi_20160714.html"},
	}

	for _, test := range tests {
		func() {
			fp, err := os.Open(test.input)
			if err != nil {
				t.Errorf("Failed to open test list data. input:%s", test.input)
			}
			defer fp.Close()

			parser := test.c.parser.(*blogChannelParser)
			items, err := parser.parseList(fp)
			if err != nil {
				t.Errorf("Failed to parse list. error:%v", err)
			}

			if len(items) != expectMaxItemNum {
				t.Errorf("Invalid items size. len:%d", len(items))
			}
		}()
	}
}

func TestBlogChannelParserItem(t *testing.T) {
	var tests = []struct {
		c                *ChannelClient
		input            string
		expectedImageLen int
		expectedVideoLen int
	}{
		{NewTamaiBlogChannelClient(nil), "testdata/blog/item_tamai_20160712.html", 6, 0},
		{NewMomotaBlogChannelClient(nil), "testdata/blog/item_momota_20160712.html", 3, 0},
		{NewAriyasuBlogChannelClient(nil), "testdata/blog/item_ariyasu_20160702.html", 0, 1},
		{NewSasakiBlogChannelClient(nil), "testdata/blog/item_sasaki_20160712.html", 2, 0},
		{NewTakagiBlogChannelClient(nil), "testdata/blog/item_takagi_20160712.html", 5, 0},
	}

	for _, test := range tests {
		func() {
			fp, err := os.Open(test.input)
			if err != nil {
				t.Errorf("Failed to open test item data. input:%s", test.input)
			}
			defer fp.Close()

			item := ChannelItem{}
			parser := test.c.parser.(*blogChannelParser)
			err = parser.parseItem(fp, &item)
			if err != nil {
				t.Errorf("Failed to parse item. error:%v", err)
			}

			if len(item.Images) != test.expectedImageLen {
				t.Errorf("Invalid image length. %d = %d", len(item.Images), test.expectedImageLen)
			}

			if len(item.Videos) != test.expectedVideoLen {
				t.Errorf("Invalid video length. %d = %d", len(item.Videos), test.expectedVideoLen)
			}
		}()
	}
}
