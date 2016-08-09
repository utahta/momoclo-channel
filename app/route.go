package momoclo_channel

import "net/http"

func init() {
	http.HandleFunc("/cron/crawl", appHandler(crawlHandler))
	http.HandleFunc("/queue/tweet", appHandler(tweetHandler))
}
