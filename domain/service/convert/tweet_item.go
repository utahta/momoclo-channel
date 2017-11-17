package convert

import (
	"strings"

	"github.com/utahta/momoclo-channel/domain/model"
)

// FeedItemToTweetItem returns TweetItem given FeedItem
func FeedItemToTweetItem(item model.FeedItem) *model.TweetItem {
	ti := &model.TweetItem{
		Title:       item.EntryTitle,
		URL:         item.EntryURL,
		PublishedAt: item.PublishedAt,
		ImageURLs:   strings.Join(item.ImageURLs, ","),
		VideoURLs:   strings.Join(item.VideoURLs, ","),
	}
	ti.BuildID()
	return ti
}
