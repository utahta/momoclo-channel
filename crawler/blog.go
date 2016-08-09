package crawler

import (
	"time"
	"regexp"
	"strings"
	"net/url"
	"fmt"
	"io"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

type BlogChannel struct {
	*Channel
}

func NewTamaiBlogChannel() *BlogChannel {
	return &BlogChannel{Channel: &Channel{Url: "http://ameblo.jp/tamai-sd/entrylist.html"}}
}

func NewMomotaBlogChannel() *BlogChannel {
	return &BlogChannel{Channel: &Channel{Url: "http://ameblo.jp/momota-sd/entrylist.html"}}
}

func NewAriyasuBlogChannel() *BlogChannel {
	return &BlogChannel{Channel: &Channel{Url: "http://ameblo.jp/ariyasu-sd/entrylist.html"}}
}

func NewSasakiBlogChannel() *BlogChannel {
	return &BlogChannel{Channel: &Channel{Url: "http://ameblo.jp/sasaki-sd/entrylist.html"}}
}

func NewTakagiBlogChannel() *BlogChannel {
	return &BlogChannel{Channel: &Channel{Url: "http://ameblo.jp/takagi-sd/entrylist.html"}}
}

func FetchTamaiBlog() ([]*ChannelItem, error) {
	return FetchParse(NewTamaiBlogChannel())
}

func FetchMomotaBlog() ([]*ChannelItem, error) {
	return FetchParse(NewMomotaBlogChannel())
}

func FetchAriyasuBlog() ([]*ChannelItem, error) {
	return FetchParse(NewAriyasuBlogChannel())
}

func FetchSasakiBlog() ([]*ChannelItem, error) {
	return FetchParse(NewSasakiBlogChannel())
}

func FetchTakagiBlog() ([]*ChannelItem, error) {
	return FetchParse(NewTakagiBlogChannel())
}

func (c *BlogChannel) Fetch() (io.ReadCloser, error) {
	return c.fetch(c.Url)
}

func (c *BlogChannel) Parse(r io.Reader) ([]*ChannelItem, error) {
	items, err := c.parseList(r)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		err := func () error {
			r, err := c.fetch(item.Url)
			if err != nil {
				return err
			}
			defer r.Close()

			err = c.parseItem(r, item)
			if err != nil {
				return err
			}
			return nil
		}()
		if err != nil {
			return nil, err
		}
	}
	return items, nil
}

func (c *BlogChannel) parseList(r io.Reader) ([]*ChannelItem, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to new document url:%s", c.Url)
	}

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to load location Asia/Tokyo")
	}

	items := []*ChannelItem{}
	err = nil
	doc.Find("[amb-component='archiveList'] > li").EachWithBreak(func(i int, s *goquery.Selection) bool {
		title := strings.TrimSpace(s.Find("[amb-component='entryItemTitle']").Text())
		href, exists := s.Find("[amb-component='entryItemTitle'] > a").Attr("href")
		if !exists {
			err = errors.Errorf("Failed to get href attribute. url:%s", c.Url)
			return false
		}

		publishedAt, err := time.ParseInLocation(
			"2006-01-02 15:04:05",
			strings.TrimSpace(strings.Replace(s.Find("[amb-component='entryItemDatetime']").Text(), "NEW!", "", 1)),
			loc,
		)
		if err != nil {
			err = errors.Wrap(err, "Failed to parse blog publish date.")
			return false
		}

		item := ChannelItem{
			Title: title,
			Url: href,
			PublishedAt: &publishedAt,
		}
		items = append(items, &item)
		return true
	})
	return items, err
}

func (c *BlogChannel) parseItem(r io.Reader, item *ChannelItem) error {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return errors.Wrapf(err, "Failed to new document. url:%s", item.Url)
	}
	if item.Images, err = c.parseImages(doc); err != nil {
		return err
	}
	if item.Videos, err = c.parseVideos(doc); err != nil {
		return err
	}
	return nil
}

func (c *BlogChannel) parseImages(doc *goquery.Document) (images []*ChannelImage, err error) {
	doc.Find("[amb-component='entryBody'] img").EachWithBreak(func(i int, s *goquery.Selection) bool {
		src, exists := s.Attr("src")
		if !exists {
			return true
		}

		matched, err := regexp.MatchString("^http://stat.ameba.jp/.*", src)
		if err != nil {
			err = errors.Wrapf(err, "Failed to regexp match string. src:%s", src)
			return false
		}
		if matched {
			images = append(images, &ChannelImage{Url: src})
		}
		return true
	})
	return
}

func (c *BlogChannel) parseVideos(doc *goquery.Document) (videos []*ChannelVideo, err error) {
	doc.Find("[amb-component='entryBody'] iframe").EachWithBreak(func(i int, s *goquery.Selection) bool {
		src, exists := s.Attr("src")
		if !exists {
			return true
		}

		matched, err := regexp.MatchString("^http://static.blog-video.jp/.*", src)
		if err != nil {
			err = errors.Wrapf(err, "Failed to regexp match string. src:%s", src)
			return false
		}
		if !matched {
			return true
		}

		u, err := url.Parse(src)
		if err != nil {
			err = errors.Wrapf(err, "Failed to parse url. src:%s", src)
			return false
		}
		src = fmt.Sprintf("http://static.blog-video.jp/output/hq/%s.mp4", u.Query().Get("v"))

		videos = append(videos, &ChannelVideo{Url: src})
		return true
	})
	return
}
