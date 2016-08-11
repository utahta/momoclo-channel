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

type blogChannelParser struct {
	channel *Channel
}

func newBlogChannelClient(url string) *ChannelClient {
	c := newChannel(url)
	return newChannelClient(c, &blogChannelParser{ channel: c })
}

func NewTamaiBlogChannelClient() *ChannelClient {
	return newBlogChannelClient("http://ameblo.jp/tamai-sd/entrylist.html")
}

func NewMomotaBlogChannelClient() *ChannelClient {
	return newBlogChannelClient("http://ameblo.jp/momota-sd/entrylist.html")
}

func NewAriyasuBlogChannelClient() *ChannelClient {
	return newBlogChannelClient("http://ameblo.jp/ariyasu-sd/entrylist.html")
}

func NewSasakiBlogChannelClient() *ChannelClient {
	return newBlogChannelClient("http://ameblo.jp/sasaki-sd/entrylist.html")
}

func NewTakagiBlogChannelClient() *ChannelClient {
	return newBlogChannelClient("http://ameblo.jp/takagi-sd/entrylist.html")
}

func FetchTamaiBlog() ([]*ChannelItem, error) {
	return NewTamaiBlogChannelClient().Fetch()
}

func FetchMomotaBlog() ([]*ChannelItem, error) {
	return NewMomotaBlogChannelClient().Fetch()
}

func FetchAriyasuBlog() ([]*ChannelItem, error) {
	return NewAriyasuBlogChannelClient().Fetch()
}

func FetchSasakiBlog() ([]*ChannelItem, error) {
	return NewSasakiBlogChannelClient().Fetch()
}

func FetchTakagiBlog() ([]*ChannelItem, error) {
	return NewTakagiBlogChannelClient().Fetch()
}

func (p *blogChannelParser) Parse(r io.Reader) ([]*ChannelItem, error) {
	c := p.channel
	items, err := p.parseList(r)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		err := func () error {
			resp, err := c.Client.Get(item.Url)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			err = p.parseItem(r, item)
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

func (p *blogChannelParser) parseList(r io.Reader) ([]*ChannelItem, error) {
	c := p.channel
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

func (p *blogChannelParser) parseItem(r io.Reader, item *ChannelItem) error {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return errors.Wrapf(err, "Failed to new document. url:%s", item.Url)
	}
	if item.Images, err = p.parseImages(doc); err != nil {
		return err
	}
	if item.Videos, err = p.parseVideos(doc); err != nil {
		return err
	}
	return nil
}

func (p *blogChannelParser) parseImages(doc *goquery.Document) (images []*ChannelImage, err error) {
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

func (p *blogChannelParser) parseVideos(doc *goquery.Document) (videos []*ChannelVideo, err error) {
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
