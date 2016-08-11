package momoclo_channel

import "net/http"

func init() {
	http.HandleFunc("/cron/crawl", appHandlerFunc(crawlHandler))
	http.HandleFunc("/queue/tweet", appHandlerFunc(tweetHandler))
}
