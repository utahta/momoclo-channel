package crawler

import (
	"testing"
	"os"
)

func TestBlogChannelParserList(t *testing.T) {
	var tests = []struct {
		c *BlogChannel
		input string
	}{
		{ NewTamaiBlogChannel(), "testdata/blog/list_tamai_20160714.html" },
		{ NewMomotaBlogChannel(), "testdata/blog/list_momota_20160714.html" },
		{ NewAriyasuBlogChannel(), "testdata/blog/list_ariyasu_20160714.html" },
		{ NewSasakiBlogChannel(), "testdata/blog/list_sasaki_20160714.html" },
		{ NewTakagiBlogChannel(), "testdata/blog/list_takagi_20160714.html" },
	}

	for _, test := range tests {
		func () {
			fp, err := os.Open(test.input)
			if err != nil {
				t.Errorf("Failed to open test list data. input:%s", test.input)
			}
			defer fp.Close()

			items, err := test.c.parseList(fp)
			if err != nil {
				t.Errorf("Failed to parse list. error:%v", err)
			}

			if len(items) != 20 {
				t.Errorf("Invalid items size. len:%d", len(items))
			}
		}()
	}
}

func TestBlogChannelParseItem(t *testing.T) {
	var tests = []struct {
		c *BlogChannel
		input string
		expectedImageLen int
		expectedVideoLen int
	}{
		{ NewTamaiBlogChannel(), "testdata/blog/item_tamai_20160712.html", 6, 0 },
		{ NewMomotaBlogChannel(), "testdata/blog/item_momota_20160712.html", 3, 0 },
		{ NewAriyasuBlogChannel(), "testdata/blog/item_ariyasu_20160702.html", 0, 1 },
		{ NewSasakiBlogChannel(), "testdata/blog/item_sasaki_20160712.html", 2, 0 },
		{ NewTakagiBlogChannel(), "testdata/blog/item_takagi_20160712.html", 5, 0 },
	}

	for _, test := range tests {
		func () {
			fp, err := os.Open(test.input)
			if err != nil {
				t.Errorf("Failed to open test item data. input:%s", test.input)
			}
			defer fp.Close()

			item := ChannelItem{}
			err = test.c.parseItem(fp, &item)
			if err != nil {
				t.Errorf("Failed to parse item. error:%v", err)
			}

			if len(item.Images) != test.expectedImageLen {
				t.Errorf("Invalid image length. %d = %d", len(item.Images), test.expectedImageLen)
			}

			if len(item.Videos) != test.expectedVideoLen{
				t.Errorf("Invalid video length. %d = %d", len(item.Videos), test.expectedVideoLen)
			}
		}()
	}
}
