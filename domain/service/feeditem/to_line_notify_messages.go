package feeditem

import (
	"fmt"

	"github.com/utahta/momoclo-channel/types"
)

// ToLineNotifyMessages converts FeedItem to []LineNotifyMessage
func ToLineNotifyMessages(item types.FeedItem) []types.LineNotifyMessage {
	var messages []types.LineNotifyMessage

	text := fmt.Sprintf("\n%s\n%s\n%s", item.Title, item.EntryTitle, item.EntryURL)
	if len(item.ImageURLs) > 0 {
		messages = append(messages, types.LineNotifyMessage{Text: text, ImageURL: item.ImageURLs[0]})
		item.ImageURLs = item.ImageURLs[1:]
	} else {
		messages = append(messages, types.LineNotifyMessage{Text: text})
	}

	for _, imageURL := range item.ImageURLs {
		messages = append(messages, types.LineNotifyMessage{Text: " ", ImageURL: imageURL}) // need space
	}

	return messages
}
