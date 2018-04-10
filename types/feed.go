package types

import (
	"fmt"
	"time"
)

const (
	FeedCodeMomota   FeedCode = "momota-sd"
	FeedCodeAriyasu  FeedCode = "ariyasu-sd"
	FeedCodeTamai    FeedCode = "tamai-sd"
	FeedCodeSasaki   FeedCode = "sasaki-sd"
	FeedCodeTakagi   FeedCode = "takagi-sd"
	FeedCodeHappyclo FeedCode = "happyclo"
	FeedCodeAeNews   FeedCode = "aenews"
	FeedCodeYoutube  FeedCode = "youtube"
)

type (
	// FeedCode represents identify code of feed
	FeedCode string

	// FeedItem represents an entry in the feed
	FeedItem struct {
		Title       string `validate:"required"`
		URL         string `validate:"required,url"`
		EntryTitle  string `validate:"required"`
		EntryURL    string `validate:"required,url"`
		ImageURLs   []string
		VideoURLs   []string
		PublishedAt time.Time `validate:"required"`
	}

	// FeedFetcher interface
	FeedFetcher interface {
		Fetch(code FeedCode, maxItemNum int, latestURL string) ([]FeedItem, error)
	}
)

// String returns string representation of FeedCode
func (f FeedCode) String() string {
	return string(f)
}

// UniqueURL builds unique url
func (i FeedItem) UniqueURL() string {
	id := i.EntryURL
	if !i.PublishedAt.IsZero() {
		id = fmt.Sprintf("%s&t=%s", id, i.PublishedAt.Format("20060102150405"))
	}
	return id
}

// ToLineNotifyMessages converts FeedItem to []LineNotifyMessage
func (i FeedItem) ToLineNotifyMessages() []LineNotifyMessage {
	var messages []LineNotifyMessage

	text := fmt.Sprintf("\n%s\n%s\n%s", i.Title, i.EntryTitle, i.EntryURL)
	if len(i.ImageURLs) > 0 {
		messages = append(messages, LineNotifyMessage{Text: text, ImageURL: i.ImageURLs[0]})
		i.ImageURLs = i.ImageURLs[1:]
	} else {
		messages = append(messages, LineNotifyMessage{Text: text})
	}

	for _, imageURL := range i.ImageURLs {
		messages = append(messages, LineNotifyMessage{Text: " ", ImageURL: imageURL}) // need space
	}
	return messages
}

// ToTweetRequests converts FeedItem to []TweetRequest
func (i FeedItem) ToTweetRequests() []TweetRequest {
	var requests []TweetRequest

	const maxUploadMediaLen = 4
	var imagesURLs [][]string
	var tmp []string
	for _, imageURL := range i.ImageURLs {
		tmp = append(tmp, imageURL)
		if len(tmp) == maxUploadMediaLen {
			imagesURLs = append(imagesURLs, tmp)
			tmp = nil
		}
	}
	if len(tmp) > 0 {
		imagesURLs = append(imagesURLs, tmp)
	}
	text := i.toTweetText()
	videoURLs := i.VideoURLs

	if len(imagesURLs) > 0 {
		requests = append(requests, TweetRequest{Text: text, ImageURLs: imagesURLs[0]})
		imagesURLs = imagesURLs[1:]
	} else if len(videoURLs) > 0 {
		requests = append(requests, TweetRequest{Text: text, VideoURL: videoURLs[0]})
		videoURLs = videoURLs[1:]
	} else {
		requests = append(requests, TweetRequest{Text: text})
	}

	if len(imagesURLs) > 0 {
		for _, imageURLs := range imagesURLs {
			requests = append(requests, TweetRequest{ImageURLs: imageURLs})
		}
	}

	if len(videoURLs) > 0 {
		for _, videoURL := range videoURLs {
			requests = append(requests, TweetRequest{VideoURL: videoURL})
		}
	}
	return requests
}

func (i FeedItem) toTweetText() string {
	const maxCharCount = 77 // max character count without hashtag and any urls TODO: correct?

	runes := []rune(fmt.Sprintf("%s %s", i.Title, i.EntryTitle))
	if len(runes) >= maxCharCount {
		runes = append(runes[0:maxCharCount-3], []rune("...")...)
	}
	return fmt.Sprintf("%s %s #momoclo #ももクロ", string(runes), i.EntryURL)
}