package crawler

import (
	"time"
	"strings"
	"fmt"
	"net/url"
	"io"

	"github.com/pkg/errors"
	"github.com/PuerkitoBio/goquery"
)

type aeNewsChannelParser struct {
	context *ChannelContext
}

func NewAeNewsChannel() *Channel {
	ctx := &ChannelContext{ Url: "http://www.momoclo.net/news/" }
	return &Channel{ Context: ctx, parser: &aeNewsChannelParser{ context: ctx } }
}

func FetchAeNews() ([]*ChannelItem, error) {
	return NewAeNewsChannel().Fetch()
}

func (p *aeNewsChannelParser) Parse(r io.Reader) ([]*ChannelItem, error) {
	ctx := p.context
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to new document. url:%s", ctx.Url)
	}

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to load location. Asia/Tokyo")
	}

	items := []*ChannelItem{}
	err = nil
	doc.Find("[class='schedule'] > [class='article']").EachWithBreak(func(i int, s *goquery.Selection) bool {
		publishedAt, err := time.ParseInLocation(
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

		u, err := url.Parse(ctx.Url)
		if err != nil {
			err = errors.Wrapf(err, "Failed to parse url. url:%s", ctx.Url)
			return false
		}

		item := ChannelItem{
			Title: a.Text(),
			Url: fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, path),
			PublishedAt: &publishedAt,
		}
		items = append(items, &item)
		return true
	})
	return items, err
}
