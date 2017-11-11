package crawler

import (
	"testing"
	"time"

	"github.com/utahta/momoclo-channel/lib/timeutil"
	"google.golang.org/appengine/aetest"
)

func Test_crawlChannelClients(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	tests := []struct {
		date      string
		expectNum int
	}{
		{"2017-01-21 00:00:00 +0900", 7},
		{"2017-01-21 07:30:00 +0900", 7},
		{"2017-01-21 16:55:00 +0900", 7},
		{"2017-01-21 20:00:00 +0900", 8},
		{"2017-01-21 20:59:00 +0900", 7},
		{"2017-01-22 16:54:00 +0900", 7},
		{"2017-01-22 16:55:00 +0900", 8},
		{"2017-01-22 16:59:00 +0900", 8},
		{"2017-01-22 19:59:00 +0900", 7},
		{"2017-01-22 20:00:00 +0900", 8},
		{"2017-01-22 20:30:00 +0900", 8},
		{"2017-01-22 23:30:00 +0900", 8},
	}

	for _, test := range tests {
		timeutil.Now = func() time.Time {
			t, _ := time.Parse("2006-01-02 15:04:05 -0700", test.date)
			return t
		}
		clients := crawlChannelClients(ctx)
		if len(clients) != test.expectNum {
			t.Errorf("Expected number of clients %d, got %d. date:%v", test.expectNum, len(clients), test.date)
		}
	}
}
