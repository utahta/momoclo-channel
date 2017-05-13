package crawler

import (
	"testing"
	"time"

	"google.golang.org/appengine/aetest"
)

func Test_crawlChannelClients(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Error(err)
	}
	defer done()

	tests := []struct {
		date      string
		expectNum int
	}{
		{"2017-01-21 16:55:00 +0900", 7},
		{"2017-01-21 17:59:00 +0900", 7},
		{"2017-01-22 16:54:00 +0900", 7},
		{"2017-01-22 16:55:00 +0900", 8},
		{"2017-01-22 17:59:00 +0900", 8},
		{"2017-01-22 18:59:00 +0900", 8},
		{"2017-01-22 19:00:00 +0900", 7},
	}

	for _, test := range tests {
		timeNow = func() time.Time {
			t, _ := time.Parse("2006-01-02 15:04:05 -0700", test.date)
			return t
		}
		clients := crawlChannelClients(ctx)
		if len(clients) != test.expectNum {
			t.Errorf("Expected number of clients %d, got %d", test.expectNum, len(clients))
		}
	}
}
