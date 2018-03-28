package feeditem

import (
	"fmt"

	"github.com/utahta/momoclo-channel/types"
)

// ToTweetRequests converts FeedItem to []TweetRequest
func ToTweetRequests(item types.FeedItem) []types.TweetRequest {
	var requests []types.TweetRequest

	const maxUploadMediaLen = 4
	var imagesURLs [][]string
	var tmp []string
	for _, imageURL := range item.ImageURLs {
		tmp = append(tmp, imageURL)
		if len(tmp) == maxUploadMediaLen {
			imagesURLs = append(imagesURLs, tmp)
			tmp = nil
		}
	}
	if len(tmp) > 0 {
		imagesURLs = append(imagesURLs, tmp)
	}
	text := toTweetText(item)
	videoURLs := item.VideoURLs

	if len(imagesURLs) > 0 {
		requests = append(requests, types.TweetRequest{Text: text, ImageURLs: imagesURLs[0]})
		imagesURLs = imagesURLs[1:]
	} else if len(videoURLs) > 0 {
		requests = append(requests, types.TweetRequest{Text: text, VideoURL: videoURLs[0]})
		videoURLs = videoURLs[1:]
	} else {
		requests = append(requests, types.TweetRequest{Text: text})
	}

	if len(imagesURLs) > 0 {
		for _, imageURLs := range imagesURLs {
			requests = append(requests, types.TweetRequest{ImageURLs: imageURLs})
		}
	}

	if len(videoURLs) > 0 {
		for _, videoURL := range videoURLs {
			requests = append(requests, types.TweetRequest{VideoURL: videoURL})
		}
	}
	return requests
}

func toTweetText(item types.FeedItem) string {
	const maxCharCount = 77 // max character count without hashtag and any urls TODO: correct?

	runes := []rune(fmt.Sprintf("%s %s", item.Title, item.EntryTitle))
	if len(runes) >= maxCharCount {
		runes = append(runes[0:maxCharCount-3], []rune("...")...)
	}
	return fmt.Sprintf("%s %s #momoclo #ももクロ", string(runes), item.EntryURL)
}
