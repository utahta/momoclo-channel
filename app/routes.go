package momoclo_channel

import (
	"net/http"
)

func init() {
	http.Handle("/cron/", &CronHandler{})
	http.Handle("/queue/", &QueueHandler{})
}
