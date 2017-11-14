package tweet

import (
	"fmt"

	"github.com/utahta/momoclo-channel/domain/model"
)

// ConvertFeedTweets converts FeedItem to FeedTweet
func ConvertFeedTweets(item model.FeedItem) []model.FeedTweet {
	var tweets []model.FeedTweet

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
	text := truncateText(item.Title, item.EntryTitle, item.EntryURL)
	videoURLs := item.VideoURLs

	if len(imagesURLs) > 0 {
		tweets = append(tweets, model.FeedTweet{Text: text, ImageURLs: imagesURLs[0]})
		imagesURLs = imagesURLs[1:]
	} else if len(videoURLs) > 0 {
		tweets = append(tweets, model.FeedTweet{Text: text, VideoURL: videoURLs[0]})
		videoURLs = videoURLs[1:]
	} else {
		tweets = append(tweets, model.FeedTweet{Text: text})
	}

	if len(imagesURLs) > 0 {
		for _, imageURLs := range imagesURLs {
			tweets = append(tweets, model.FeedTweet{ImageURLs: imageURLs})
		}
	}

	if len(videoURLs) > 0 {
		for _, videoURL := range videoURLs {
			tweets = append(tweets, model.FeedTweet{VideoURL: videoURL})
		}
	}
	return tweets
}

func truncateText(title, entryTitle, entryURL string) string {
	const maxTextLen = 77 // max text length without hashtag and image url

	runes := []rune(fmt.Sprintf("%s %s", title, entryTitle))
	if len(runes) >= maxTextLen {
		runes = append(runes[0:maxTextLen-3], []rune("...")...)
	}
	return fmt.Sprintf("%s %s #momoclo #ももクロ", string(runes), entryURL)
}
