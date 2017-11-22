package feeditem

import (
	"fmt"
	"strings"

	"github.com/utahta/momoclo-channel/domain/model"
)

// ToLineItem returns LineItem given FeedItem
func ToLineItem(item model.FeedItem) *model.LineItem {
	ti := &model.LineItem{
		Title:       item.EntryTitle,
		URL:         item.EntryURL,
		PublishedAt: item.PublishedAt,
		ImageURLs:   strings.Join(item.ImageURLs, ","),
		VideoURLs:   strings.Join(item.VideoURLs, ","),
	}
	ti.BuildID()
	return ti
}

// ToLineNotifyRequests converts FeedItem to []LineNotifyRequest
func ToLineNotifyRequests(item model.FeedItem) []model.LineNotifyRequest {
	var requests []model.LineNotifyRequest

	text := fmt.Sprintf("\n%s\n%s\n%s", item.Title, item.EntryTitle, item.EntryURL)
	if len(item.ImageURLs) > 0 {
		requests = append(requests, model.LineNotifyRequest{Text: text, ImageURL: item.ImageURLs[0]})
		item.ImageURLs = item.ImageURLs[1:]
	} else {
		requests = append(requests, model.LineNotifyRequest{Text: text})
	}

	for _, imageURL := range item.ImageURLs {
		requests = append(requests, model.LineNotifyRequest{ImageURL: imageURL})
	}

	return requests
}
