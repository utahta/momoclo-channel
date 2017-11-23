package feeditem

import (
	"fmt"

	"github.com/utahta/momoclo-channel/domain/model"
)

// ToLineNotifyMessages converts FeedItem to []LineNotifyMessage
func ToLineNotifyMessages(item model.FeedItem) []model.LineNotifyMessage {
	var messages []model.LineNotifyMessage

	text := fmt.Sprintf("\n%s\n%s\n%s", item.Title, item.EntryTitle, item.EntryURL)
	if len(item.ImageURLs) > 0 {
		messages = append(messages, model.LineNotifyMessage{Text: text, ImageURL: item.ImageURLs[0]})
		item.ImageURLs = item.ImageURLs[1:]
	} else {
		messages = append(messages, model.LineNotifyMessage{Text: text})
	}

	for _, imageURL := range item.ImageURLs {
		messages = append(messages, model.LineNotifyMessage{Text: " ", ImageURL: imageURL}) // need space
	}

	return messages
}
