package crawler

import (
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

type aeNewsChannelParser struct {
	channel *Channel
}

func NewAeNewsChannelClient(options ...ChannelClientOption) (*ChannelClient, error) {
	c := newChannel("http://www.momoclo.net/news/", "ANGEL EYES | News")
	return newChannelClient(c, &aeNewsChannelParser{channel: c}, options...)
}

func FetchAeNews() (*Channel, error) {
	cc, err := NewAeNewsChannelClient()
	if err != nil {
		return nil, err
	}
	return cc.Fetch()
}

func (p *aeNewsChannelParser) Parse(r io.Reader) ([]*ChannelItem, error) {
	c := p.channel
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to new document. url:%s", c.Url)
	}

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to load location. Asia/Tokyo")
	}

	items := []*ChannelItem{}
	err = nil
	doc.Find("[class='schedule'] > [class='article']").EachWithBreak(func(i int, s *goquery.Selection) bool {
		var publishedAt time.Time
		publishedAt, err = time.ParseInLocation(
			"2006/01/02",
			strings.TrimSpace(fmt.Sprintf("%s/%s", s.Find("[class='year clearfix']").Text(), s.Find("[class='date clearfix']").Text())),
			loc,
		)
		if err != nil {
			err = errors.Wrapf(err, "Failed to parse in location. i:%d", i)
			return false
		}

		a := s.Find("h4 > a").First()
		path, exists := a.Attr("href")
		if !exists {
			err = errors.Errorf("Failed to get href attribute. a:%v", a)
			return false
		}

		var u *url.URL
		u, err = url.Parse(c.Url)
		if err != nil {
			err = errors.Wrapf(err, "Failed to parse url. url:%s", c.Url)
			return false
		}

		item := ChannelItem{
			Title:       a.Text(),
			Url:         fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, path),
			PublishedAt: &publishedAt,
		}
		items = append(items, &item)
		return true
	})
	return items, err
}
