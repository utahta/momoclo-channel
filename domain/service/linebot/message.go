package linebot

import (
	"fmt"

	"github.com/utahta/momoclo-channel/lib/config"
)

// FollowMessage returns message on follow
func FollowMessage() string {
	return fmt.Sprintf(`友だち追加ありがとうございます。
こちらは、ももクロちゃんのブログやAE NEWS等を通知する機能との連携を補助したり、画像を返したりするBOTです。

%s

%s
`, HelpMessage(), OnMessage())
}

// HelpMessage returns help message
func HelpMessage() string {
	urlStr := fmt.Sprintf("%s%s", config.C.App.BaseURL, "/line/bot/help")
	return fmt.Sprintf("ヘルプ（・Θ・）\n%s", urlStr)
}

// OnMessage returns line notification on message
func OnMessage() string {
	urlStr := fmt.Sprintf("%s%s", config.C.App.BaseURL, "/line/notify/on")
	return fmt.Sprintf("通知機能を有効にする場合は、下記URLをクリックしてください（・Θ・）\n%s", urlStr)
}

// OffMessage returns line notification off message
func OffMessage() string {
	return "通知機能を無効にする場合は、下記URLから解除を行ってください（・Θ・）\nhttps://notify-bot.line.me/my/"
}

// ImageNotFoundMessage returns image not found message
func ImageNotFoundMessage() string {
	return "画像がみつかりませんでした（・Θ・）"
}
